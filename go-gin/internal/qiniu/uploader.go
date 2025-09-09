package qiniu

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	qbox "github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Config struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string // e.g. https://cdn.example.com
}

type Uploader struct {
	cfg       Config
	mac       *qbox.Mac
	putPolicy storage.PutPolicy
	form      *storage.FormUploader
	resume    *storage.ResumeUploaderV2
}

func NewUploader(cfg Config) *Uploader {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	upCfg := &storage.Config{UseHTTPS: true, UseCdnDomains: true}
	form := storage.NewFormUploader(upCfg)
	resume := storage.NewResumeUploaderV2(upCfg)
	return &Uploader{
		cfg:       cfg,
		mac:       mac,
		putPolicy: storage.PutPolicy{Scope: cfg.Bucket},
		form:      form,
		resume:    resume,
	}
}

// Put uploads content to Qiniu at key and returns the public URL (Domain/key).
// If size < 0, it will buffer the reader in memory to determine size.
func (u *Uploader) Put(ctx context.Context, key string, r io.Reader, size int64, mime string) (string, error) {
	log.Printf("qiniu: starting upload key=%s size=%d mime=%s", key, size, mime)
	start := time.Now()

	if u == nil || u.form == nil {
		log.Printf("qiniu: uploader not initialized")
		return "", fmt.Errorf("uploader not initialized")
	}

	uploadCtx := ctx
	if dl, ok := ctx.Deadline(); !ok || time.Until(dl) < 3*time.Minute {
		var cancel context.CancelFunc
		uploadCtx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		_ = dl
		log.Printf("qiniu: applied upload timeout 5m")
	}

	if size < 0 {
		log.Printf("qiniu: buffering data to determine size")
		bufStart := time.Now()
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			log.Printf("qiniu: buffering failed after %v: %v", time.Since(bufStart), err)
			return "", err
		}
		size = int64(buf.Len())
		r = bytes.NewReader(buf.Bytes())
		log.Printf("qiniu: buffering completed in %v, actual size=%d", time.Since(bufStart), size)
	}

	log.Printf("qiniu: generating upload token")
	token := u.putPolicy.UploadToken(u.mac)

	log.Printf("qiniu: starting actual upload to bucket=%s", u.cfg.Bucket)
	uploadStart := time.Now()

	if size > 8*1024*1024 {
		url, err := u.putResumableV2(uploadCtx, token, key, r, size, mime)
		if err != nil {
			elapsed := time.Since(uploadStart)
			log.Printf("qiniu: resumable upload failed after %v: %v", elapsed, err)
			return "", err
		}
		return url, nil
	}

	var ret storage.PutRet
	extra := &storage.PutExtra{MimeType: mime}
	if size > 0 {
		r = newProgressReader(r, size)
	}
	if err := u.form.Put(uploadCtx, &ret, token, key, r, size, extra); err != nil {
		elapsed := time.Since(uploadStart)
		log.Printf("qiniu: upload failed after %v: %v", elapsed, err)
		return "", err
	}

	uploadElapsed := time.Since(uploadStart)
	totalElapsed := time.Since(start)

	domain := strings.TrimRight(u.cfg.Domain, "/")
	publicURL := ret.Key
	if domain != "" {
		publicURL = domain + "/" + ret.Key
	}

	log.Printf("qiniu: upload completed in %v (total: %v), key=%s, url=%s", uploadElapsed, totalElapsed, ret.Key, publicURL)

	return publicURL, nil
}

// Fetch pulls a remote URL directly into the bucket at key and returns the public URL.
// This avoids uploading from local network and is usually faster and more reliable.
func (u *Uploader) Fetch(ctx context.Context, key string, remoteURL string) (string, error) {
	log.Printf("qiniu: starting server-side fetch key=%s src=%s", key, remoteURL)
	start := time.Now()

	if u == nil || u.mac == nil {
		return "", fmt.Errorf("uploader not initialized")
	}

	upCfg := &storage.Config{UseHTTPS: true, UseCdnDomains: true}
	bm := storage.NewBucketManager(u.mac, upCfg)
	ret, err := bm.Fetch(remoteURL, u.cfg.Bucket, key)
	if err != nil {
		log.Printf("qiniu: fetch failed after %v: %v", time.Since(start), err)
		return "", err
	}

	domain := strings.TrimRight(u.cfg.Domain, "/")
	publicURL := ret.Key
	if domain != "" {
		publicURL = domain + "/" + ret.Key
	}
	log.Printf("qiniu: fetch completed in %v, key=%s, url=%s", time.Since(start), ret.Key, publicURL)
	return publicURL, nil
}

func (u *Uploader) putResumableV2(ctx context.Context, uptoken, key string, r io.Reader, size int64, mime string) (string, error) {
	if u.resume == nil {
		return "", fmt.Errorf("resume uploader not initialized")
	}
	tmp, err := os.CreateTemp("", "qiniu-up-*.tmp")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	log.Printf("qiniu: buffering to temp file for resumable upload: %s", tmpPath)
	bufStart := time.Now()
	src := r
	if size > 0 {
		src = newProgressReader(r, size)
	}
	n, err := io.Copy(tmp, src)
	_ = tmp.Close()
	if err != nil {
		log.Printf("qiniu: buffering to temp failed after %v: %v", time.Since(bufStart), err)
		return "", err
	}
	log.Printf("qiniu: buffering to temp completed in %v, bytes=%d", time.Since(bufStart), n)
	if size > 0 && n != size {
		log.Printf("qiniu: warning, buffered size mismatch: declared=%d copied=%d", size, n)
		size = n
	}

	file, err := os.Open(tmpPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var ret storage.PutRet
	fr, _ := storage.NewFileRecorder(os.TempDir())
	extra := &storage.RputV2Extra{PartSize: 4 * 1024 * 1024, MimeType: mime, Recorder: fr, TryTimes: 3}
	upStart := time.Now()
	if err := u.resume.Put(ctx, &ret, uptoken, key, file, size, extra); err != nil {
		log.Printf("qiniu: resumable v2 put failed after %v: %v", time.Since(upStart), err)
		return "", err
	}
	log.Printf("qiniu: resumable v2 upload ok in %v, key=%s", time.Since(upStart), ret.Key)

	domain := strings.TrimRight(u.cfg.Domain, "/")
	publicURL := ret.Key
	if domain != "" {
		publicURL = domain + "/" + ret.Key
	}
	return publicURL, nil
}

type progressReader struct {
	r         io.Reader
	total     int64
	read      int64
	nextBytes int64
	start     time.Time
	lastTick  time.Time
	stepBytes int64
	stepTick  time.Duration
}

func newProgressReader(r io.Reader, total int64) *progressReader {
	step := int64(2 * 1024 * 1024)
	return &progressReader{r: r, total: total, nextBytes: step, start: time.Now(), lastTick: time.Now(), stepBytes: step, stepTick: 2 * time.Second}
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	if n > 0 {
		newRead := atomic.AddInt64(&p.read, int64(n))
		now := time.Now()
		if newRead >= p.nextBytes || now.Sub(p.lastTick) >= p.stepTick {
			percent := float64(newRead) * 100 / float64(p.total)
			elapsed := now.Sub(p.start).Seconds()
			speedKBs := float64(newRead) / 1024 / elapsed
			log.Printf("qiniu: upload progress %.1f%% (%d/%d bytes, %.1f KB/s)", percent, newRead, p.total, speedKBs)
			p.lastTick = now
			for p.nextBytes <= newRead {
				p.nextBytes += p.stepBytes
			}
		}
	}
	return n, err
}

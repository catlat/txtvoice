package qiniu

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

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
}

func NewUploader(cfg Config) *Uploader {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	upCfg := &storage.Config{UseHTTPS: true}
	form := storage.NewFormUploader(upCfg)
	return &Uploader{
		cfg:       cfg,
		mac:       mac,
		putPolicy: storage.PutPolicy{Scope: cfg.Bucket},
		form:      form,
	}
}

// Put uploads content to Qiniu at key and returns the public URL (Domain/key).
// If size < 0, it will buffer the reader in memory to determine size.
func (u *Uploader) Put(ctx context.Context, key string, r io.Reader, size int64, mime string) (string, error) {
	if u == nil || u.form == nil {
		return "", fmt.Errorf("uploader not initialized")
	}
	if size < 0 {
		// Buffer to memory (suitable for thumbnails and small files)
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return "", err
		}
		size = int64(buf.Len())
		r = bytes.NewReader(buf.Bytes())
	}
	token := u.putPolicy.UploadToken(u.mac)
	var ret storage.PutRet
	extra := &storage.PutExtra{MimeType: mime}
	if err := u.form.Put(ctx, &ret, token, key, r, size, extra); err != nil {
		return "", err
	}
	domain := strings.TrimRight(u.cfg.Domain, "/")
	if domain == "" {
		return ret.Key, nil
	}
	return domain + "/" + ret.Key, nil
}

package dlyt

import (
	"bytes"
	"context"
)

// UploadBytesToQiniu uploads raw bytes under a given key with mime type, returning the public URL.
// If uploader is not configured, returns an empty URL with an error.
func UploadBytesToQiniu(ctx context.Context, key string, data []byte, mime string) (string, error) {
	if uploader == nil {
		return "", ErrQiniuNotConfigured
	}
	return uploader.Put(ctx, key, bytes.NewReader(data), int64(len(data)), mime)
}

// FetchToQiniu server-side fetches a remote URL to the bucket under key, returning the public URL.
// If uploader is not configured, returns an empty URL with an error.
func FetchToQiniu(ctx context.Context, key string, remoteURL string) (string, error) {
	if uploader == nil {
		return "", ErrQiniuNotConfigured
	}
	return uploader.Fetch(ctx, key, remoteURL)
}

// ErrQiniuNotConfigured is returned when Qiniu uploader is not initialized.
var ErrQiniuNotConfigured = errQiniuNotConfigured{}

type errQiniuNotConfigured struct{}

func (errQiniuNotConfigured) Error() string { return "qiniu uploader not configured" }

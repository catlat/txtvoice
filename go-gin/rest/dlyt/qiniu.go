package dlyt

import "go-gin/internal/qiniu"

var uploader *qiniu.Uploader

func SetQiniuConfig(cfg qiniu.Config) {
	uploader = qiniu.NewUploader(cfg)
}

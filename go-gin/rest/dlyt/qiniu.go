package dlyt

import "go-gin/internal/qiniu"

var uploader *qiniu.Uploader
var qiniuDomain string

func SetQiniuConfig(cfg qiniu.Config) {
	uploader = qiniu.NewUploader(cfg)
	qiniuDomain = cfg.Domain
}

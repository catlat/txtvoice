package controller

import (
	"encoding/base64"
	"go-gin/internal/httpx"
	"go-gin/internal/httpx/validators"
	"go-gin/logic"
	"go-gin/rest/dlyt"
	"go-gin/typing"
	"os"
	"path/filepath"
	"strings"
)

type ytController struct{}

var YtController = &ytController{}

func (c *ytController) Info(ctx *httpx.Context) (any, error) {
	var req typing.YtInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		return nil, err
	}
	if err := validators.Validate(&req); err != nil {
		return nil, err
	}
	return dlyt.Svc.InfoWithPlatform(ctx, req.IdOrUrl, req.Platform)
}

func (c *ytController) Text(ctx *httpx.Context) (any, error) {
	var req typing.YtTextReq
	if err := ctx.ShouldBind(&req); err != nil {
		return nil, err
	}
	if err := validators.Validate(&req); err != nil {
		return nil, err
	}
	identity := ctx.Query("identity")
	l := logic.NewTranscriptLogic()
	tr, err := l.GetOrCreateWithPlatform(ctx, req.IdOrUrl, req.TargetLan, identity, req.Platform)
	if err != nil {
		return nil, err
	}

	// 获取音频并返回原始数据（若为本地静态文件则读文件返回 data URL）
	audioData := ""
	audioType := ""
	if a, aerr := dlyt.Svc.AudioWithPlatform(ctx, req.IdOrUrl, req.Platform); aerr == nil && a != nil {
		if strings.HasPrefix(a.AudioUrl, "/static/") {
			localPath := filepath.Join("public", strings.TrimPrefix(a.AudioUrl, "/static/"))
			if b, rerr := os.ReadFile(localPath); rerr == nil {
				ext := strings.ToLower(filepath.Ext(localPath))
				mime := "audio/m4a"
				audioType = "m4a" // 默认类型
				switch ext {
				case ".mp3":
					mime = "audio/mpeg"
					audioType = "mp3"
				case ".wav":
					mime = "audio/wav"
					audioType = "wav"
				case ".m4a":
					mime = "audio/m4a"
					audioType = "m4a"
				case ".aac":
					mime = "audio/aac"
					audioType = "aac"
				case ".ogg":
					mime = "audio/ogg"
					audioType = "ogg"
				case ".flac":
					mime = "audio/flac"
					audioType = "flac"
				}
				audioData = "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(b)
			}
		}
	}

	return &typing.YtTextReply{
		OriginalText:   tr.OriginalText,
		TranslatedText: tr.TranslatedText,
		AudioData:      audioData,
		AudioType:      audioType,
	}, nil
}

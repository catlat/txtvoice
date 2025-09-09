package controller

import (
	"go-gin/internal/httpx"
	"go-gin/internal/httpx/validators"
	"go-gin/logic"
)

type adminController struct{}

var AdminController = &adminController{}

type SeedQuotaReq struct {
	Identities []string `json:"identities" binding:"required" label:"账号列表"`
	AsrChars   int      `json:"asr_chars" binding:"required" label:"ASR额度"`
	TtsChars   int      `json:"tts_chars" binding:"required" label:"TTS额度"`
}

func (c *adminController) SeedQuota(ctx *httpx.Context) (any, error) {
	var req SeedQuotaReq
	if err := ctx.ShouldBind(&req); err != nil {
		return nil, err
	}
	if err := validators.Validate(&req); err != nil {
		return nil, err
	}
	l := logic.NewAdminLogic()
	if err := l.SeedQuota(ctx, req.Identities, req.AsrChars, req.TtsChars); err != nil {
		return nil, err
	}
	return map[string]any{"ok": true}, nil
}

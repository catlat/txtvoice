package controller

import (
	"go-gin/internal/component/logx"
	"go-gin/internal/httpx"
	"go-gin/internal/httpx/validators"
	"go-gin/logic"
	"go-gin/typing"
	"strings"
)

type ttsController struct{}

var TTSController = &ttsController{}

func (c *ttsController) Synthesize(ctx *httpx.Context) (any, error) {
	var req typing.TTSSynthesizeReq
	if err := ctx.ShouldBind(&req); err != nil {
		return nil, err
	}
	if err := validators.Validate(&req); err != nil {
		return nil, err
	}
	identity := httpx.Identity(ctx)
	l := logic.NewTTSLogic()
	item, err := l.Synthesize(ctx, identity, req.Text, req.Speaker, req.UseMyVoice)
	if err != nil {
		return nil, err
	}
	resp := map[string]any{"audio_url": item.AudioUrl, "char_count": item.CharCount}
	// 业务日志：返回给前端的关键字段（避免打印巨大 data URL 全量，仅打印类型与长度）
	urlType := "remote"
	if strings.HasPrefix(item.AudioUrl, "data:") {
		urlType = "data"
	}
	logx.WithContext(ctx).Info(
		"tts_synthesize_reply",
		map[string]any{
			"audio_url_type": urlType,
			"audio_url_len":  len(item.AudioUrl),
			"char_count":     item.CharCount,
			"speaker":        req.Speaker,
			"use_my_voice":   req.UseMyVoice,
		},
	)
	return resp, nil
}

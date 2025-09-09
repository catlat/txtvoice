package controller

import (
	"go-gin/internal/httpx"
	"go-gin/internal/httpx/validators"
	"go-gin/logic"
	"go-gin/rest/dlyt"
	"go-gin/typing"
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

	return &typing.YtTextReply{
		TranslatedText: tr.TranslatedText,
	}, nil
}

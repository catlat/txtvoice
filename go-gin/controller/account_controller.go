package controller

import (
	"go-gin/internal/httpx"
	"go-gin/logic"
)

type accountController struct{}

var AccountController = &accountController{}

func (c *accountController) Profile(ctx *httpx.Context) (any, error) {
	identity := ctx.Query("identity")
	l := logic.NewAccountLogic()
	data, _ := l.Profile(ctx, identity)
	return data, nil
}

func (c *accountController) Packages(ctx *httpx.Context) (any, error) {
	identity := ctx.Query("identity")
	l := logic.NewAccountLogic()
	items, _ := l.Packages(ctx, identity)
	return map[string]any{"items": items}, nil
}

func (c *accountController) Usage(ctx *httpx.Context) (any, error) {
	identity := ctx.Query("identity")
	l := logic.NewAccountLogic()
	items, _ := l.Usage(ctx, identity)
	return map[string]any{"days": items}, nil
}

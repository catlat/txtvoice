package controller

import (
	"go-gin/internal/httpx"
	"go-gin/internal/httpx/validators"
	"go-gin/logic"
	"go-gin/typing"
	"net/http"
	"time"
)

type authController struct{}

var AuthController = &authController{}

func (c *authController) LoginSimple(ctx *httpx.Context) (any, error) {
	var req typing.LoginSimpleReq
	if err := ctx.ShouldBind(&req); err != nil {
		return nil, err
	}
	if err := validators.Validate(&req); err != nil {
		return nil, err
	}
	l := logic.NewAuthLogic()
	resp, err := l.LoginSimple(ctx, req.Identity)
	if err != nil {
		return nil, err
	}
	// set http-only cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	return resp, nil
}

func (c *authController) Logout(ctx *httpx.Context) (any, error) {
	tokenId := ctx.Query("token")
	if tokenId == "" {
		if cookie, err := ctx.Request.Cookie("token"); err == nil {
			tokenId = cookie.Value
		}
	}
	l := logic.NewAuthLogic()
	if err := l.Logout(ctx, tokenId); err != nil {
		return nil, err
	}
	// clear cookie
	http.SetCookie(ctx.Writer, &http.Cookie{Name: "token", Value: "", Path: "/", HttpOnly: true, Expires: time.Unix(0, 0)})
	return map[string]any{"logout": true}, nil
}

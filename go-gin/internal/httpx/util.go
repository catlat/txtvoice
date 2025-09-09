package httpx

import (
	"go-gin/internal/token"
)

// Identity resolves user identity from request:
// 1) query param `identity`
// 2) token from cookie/header -> lookup in token store
func Identity(ctx *Context) string {
	id := ctx.Query("identity")
	if id != "" {
		return id
	}
	tokenId := ctx.GetHeader("token")
	if tokenId == "" {
		if c, err := ctx.Request.Cookie("token"); err == nil {
			tokenId = c.Value
		}
	}
	if tokenId == "" {
		return ""
	}
	if v, err := token.Get(ctx, tokenId, "identity"); err == nil {
		return v
	}
	return ""
}


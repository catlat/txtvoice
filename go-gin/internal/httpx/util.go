package httpx

import (
	"go-gin/internal/token"
)

// Identity resolves user identity from request strictly via token
// 1) token from header/cookie -> lookup in token store
// 2) no query fallback is allowed for protected endpoints
func Identity(ctx *Context) string {
	tokenId := ctx.GetHeader("token")
	if tokenId == "" {
		if c, err := ctx.Request.Cookie("token"); err == nil {
			tokenId = c.Value
		}
	}
	if tokenId != "" {
		if v, err := token.Get(ctx, tokenId, "identity"); err == nil && v != "" {
			return v
		}
	}
	return ""
}

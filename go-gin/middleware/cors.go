package middleware

import (
	"go-gin/internal/httpx"
)

// CORS 允许跨域访问（开发阶段）
func CORS() httpx.HandlerFunc {
	return func(ctx *httpx.Context) (any, error) {
		origin := ctx.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		ctx.Header("Access-Control-Allow-Origin", origin)
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

		// 动态回显浏览器预检中声明的自定义头，追加常用头与 token 头
		reqHeaders := ctx.Request.Header.Get("Access-Control-Request-Headers")
		if reqHeaders == "" {
			reqHeaders = "Content-Type,Authorization,X-Requested-With,token"
		} else {
			reqHeaders = reqHeaders + ",Content-Type,Authorization,X-Requested-With,token"
		}
		ctx.Header("Access-Control-Allow-Headers", reqHeaders)
		// 暴露部分头（可选）
		ctx.Header("Access-Control-Expose-Headers", "Content-Length,Content-Type,Trace-Id")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return nil, nil
		}
		return nil, nil
	}
}

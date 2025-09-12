package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
	"go-gin/middleware"
)

// RegisterYtRoutes 注册 YouTube 相关路由
func RegisterYtRoutes(r *httpx.RouterGroup) {
	// 使用子组绑定 TokenCheck，避免污染父组其它路由
	g := r.Group("")
	g.Before(middleware.TokenCheck()).POST("/yt/info", controller.YtController.Info)
	g.Before(middleware.TokenCheck()).POST("/yt/text", controller.YtController.Text)
}

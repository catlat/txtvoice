package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
)

// RegisterYtRoutes 注册 YouTube 相关路由
func RegisterYtRoutes(r *httpx.RouterGroup) {
	r.POST("/yt/info", controller.YtController.Info)
	r.POST("/yt/text", controller.YtController.Text)
}


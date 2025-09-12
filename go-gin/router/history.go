package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
	"go-gin/middleware"
)

// RegisterHistoryRoutes 注册历史查询路由
func RegisterHistoryRoutes(r *httpx.RouterGroup) {
	g := r.Group("")
	g.Before(middleware.TokenCheck()).GET("/history/tts", controller.HistoryController.ListTTS)
}

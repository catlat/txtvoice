package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
	"go-gin/middleware"
)

// RegisterTTSRoutes 注册 TTS 路由
func RegisterTTSRoutes(r *httpx.RouterGroup) {
	g := r.Group("")
	g.Before(middleware.TokenCheck()).POST("/tts/synthesize", controller.TTSController.Synthesize)
}

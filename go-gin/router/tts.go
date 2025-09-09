package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
)

// RegisterTTSRoutes 注册 TTS 路由
func RegisterTTSRoutes(r *httpx.RouterGroup) {
	r.POST("/tts/synthesize", controller.TTSController.Synthesize)
}


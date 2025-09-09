package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
)

// RegisterHistoryRoutes 注册历史查询路由
func RegisterHistoryRoutes(r *httpx.RouterGroup) {
	r.GET("/history/videos", controller.HistoryController.ListVideos)
	r.GET("/history/video/:source_site/:video_id", controller.HistoryController.GetVideoDetail)
	r.GET("/history/tts", controller.HistoryController.ListTTS)
}


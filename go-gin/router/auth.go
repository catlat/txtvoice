package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
)

// RegisterAuthRoutes 注册认证路由
func RegisterAuthRoutes(r *httpx.RouterGroup) {
	r.POST("/auth/login_simple", controller.AuthController.LoginSimple)
	r.POST("/auth/logout", controller.AuthController.Logout)
}


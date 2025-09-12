package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
	"go-gin/middleware"
)

// RegisterAuthRoutes 注册认证路由
func RegisterAuthRoutes(r *httpx.RouterGroup) {
	// 新登录：手机号+密码
	r.POST("/auth/login", controller.AuthController.Login)
	// 修改密码（需要登录）
	r.Before(middleware.TokenCheck()).POST("/auth/change_password", controller.AuthController.ChangePassword)
	// 登出
	r.POST("/auth/logout", controller.AuthController.Logout)
}

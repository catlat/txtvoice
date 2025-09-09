package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
)

// RegisterAccountRoutes 注册账号相关路由
func RegisterAccountRoutes(r *httpx.RouterGroup) {
	r.GET("/account/profile", controller.AccountController.Profile)
	r.GET("/account/packages", controller.AccountController.Packages)
	r.GET("/account/usage", controller.AccountController.Usage)
}


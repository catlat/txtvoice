package router

import (
	"go-gin/controller"
	"go-gin/internal/httpx"
)

// RegisterAdminRoutes 临时管理接口（内测期间）
func RegisterAdminRoutes(r *httpx.RouterGroup) {
	r.POST("/admin/seed_quota", controller.AdminController.SeedQuota)
}

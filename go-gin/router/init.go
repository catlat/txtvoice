package router

import (
	"go-gin/internal/httpx"
)

func Init(route *httpx.Engine) {
	RegisterCommonRoutes(route)
	RegisterLoginRoutes(route.Group("/"))

	api := route.Group("/api")
	RegisterApiRoutes(api)
	RegisterYtRoutes(api)
	RegisterTTSRoutes(api)
	RegisterHistoryRoutes(api)
	RegisterAccountRoutes(api)
	RegisterAuthRoutes(api)
	RegisterAdminRoutes(api)

	RegisterDemoRoutes(route.Group("/demo"))

	// Serve local static files under /static from ./public
	route.Engine.Static("/static", "./public")
}

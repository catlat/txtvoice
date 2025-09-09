package router

import (
	"go-gin/internal/component/db"
	"go-gin/internal/component/redisx"
	"go-gin/internal/errorx"
	"go-gin/internal/g"
	"go-gin/internal/httpx"
	"go-gin/util"
)

func RegisterCommonRoutes(route *httpx.Engine) {
	route.NoMethod(func(ctx *httpx.Context) (any, error) {
		return nil, errorx.ErrMethodNotAllowed
	})
	route.NoRoute(func(ctx *httpx.Context) (any, error) {
		return nil, errorx.ErrNoRoute
	})
	// 健康检测
	route.GET("/status", func(ctx *httpx.Context) (any, error) {
		db_err := db.WithContext(ctx).Ping()
		redis_err := redisx.Client().Ping(ctx).Err()
		return g.MapStrStr{
			"database": util.When(db_err == nil, "ok", "failed"),
			"redis":    util.When(redis_err == nil, "ok", "failed"),
		}, nil
	})
	// 根路径简单返回
	route.GET("/", func(ctx *httpx.Context) (any, error) {
		return g.MapStrStr{"service": "ok"}, nil
	})
}

package controller

import (
	"go-gin/internal/component/logx"
	"go-gin/internal/httpx"
	"go-gin/logic"
	"strconv"
)

type historyController struct{}

var HistoryController = &historyController{}

func (c *historyController) ListVideos(ctx *httpx.Context) (any, error) {
	l := logic.NewHistoryLogic()
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	items, _ := l.ListVideos(ctx, page, size)
	return map[string]any{"items": items}, nil
}

func (c *historyController) GetVideoDetail(ctx *httpx.Context) (any, error) {
	l := logic.NewHistoryLogic()
	source := ctx.Param("source_site")
	vid := ctx.Param("video_id")
	detail, _ := l.GetVideoDetail(ctx, source, vid)
	return map[string]any{"detail": detail}, nil
}

func (c *historyController) ListTTS(ctx *httpx.Context) (any, error) {
	l := logic.NewHistoryLogic()
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	identity := httpx.Identity(ctx)
	items, _ := l.ListTTS(ctx, page, size, identity)
	logx.WithContext(ctx).Info("history_list_tts", map[string]any{"identity": identity, "count": len(items)})
	return map[string]any{"items": items}, nil
}

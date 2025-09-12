package controller

import (
	"go-gin/internal/component/logx"
	"go-gin/internal/httpx"
	"go-gin/logic"
	"strconv"
)

type historyController struct{}

var HistoryController = &historyController{}

// 已移除视频历史接口，仅保留合成历史

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
	items, total, _ := l.ListTTS(ctx, page, size, identity)

	// 计算分页信息
	totalPages := (total + int64(size) - 1) / int64(size) // 向上取整
	hasNext := page < int(totalPages)
	hasPrev := page > 1

	logx.WithContext(ctx).Info("history_list_tts", map[string]any{
		"identity": identity,
		"count":    len(items),
		"total":    total,
		"page":     page,
		"size":     size,
	})

	return map[string]any{
		"items": items,
		"pagination": map[string]any{
			"page":        page,
			"size":        size,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	}, nil
}

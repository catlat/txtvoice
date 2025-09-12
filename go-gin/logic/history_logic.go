package logic

import (
	"context"
	"go-gin/internal/component/db"
	"go-gin/model"
)

type HistoryLogic struct{}

func NewHistoryLogic() *HistoryLogic { return &HistoryLogic{} }

// 已移除视频历史查询逻辑，仅保留合成历史

func (l *HistoryLogic) ListTTS(ctx context.Context, page, size int, identity string) ([]model.TTSHistory, int64, error) {
	var items []model.TTSHistory
	var total int64

	// 未登录/无身份：直接返回空，避免泄露任何数据
	if identity == "" {
		return []model.TTSHistory{}, 0, nil
	}

	// 严格按用户过滤
	baseQuery := db.WithContext(ctx).Model(&model.TTSHistory{}).Where("user_identity=?", identity)

	// 获取总数
	_ = baseQuery.Count(&total)

	// 获取分页数据
	q := baseQuery.Order("id desc").Limit(size).Offset((page - 1) * size)
	_ = q.Find(&items)

	return items, total, nil
}

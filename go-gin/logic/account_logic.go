package logic

import (
	"context"
	"go-gin/internal/component/db"
	"go-gin/model"
	"time"
)

type AccountLogic struct{}

func NewAccountLogic() *AccountLogic { return &AccountLogic{} }

func (l *AccountLogic) Profile(ctx context.Context, identity string) (map[string]any, error) {
	var au model.AccountUser
	_ = db.WithContext(ctx).Where("identity=?", identity).First(&au)
	return map[string]any{
		"exists":       au.Id != 0,
		"identity":     identity,
		"status":       au.Status,
		"display_name": au.DisplayName,
	}, nil
}

func (l *AccountLogic) Packages(ctx context.Context, identity string) ([]model.UserPackage, error) {
	var items []model.UserPackage
	_ = db.WithContext(ctx).Where("user_identity=?", identity).Find(&items)
	return items, nil
}

func (l *AccountLogic) Usage(ctx context.Context, identity string) ([]model.UsageDaily, error) {
	var items []model.UsageDaily
	_ = db.WithContext(ctx).Where("user_identity=?", identity).Order("date desc").Limit(30).Find(&items)
	// 零填充近30天
	today := time.Now()
	dayToItem := map[string]model.UsageDaily{}
	for _, it := range items {
		dayToItem[it.Date] = it
	}
	var result []model.UsageDaily
	for i := 0; i < 30; i++ {
		d := today.AddDate(0, 0, -i).Format("2006-01-02")
		if v, ok := dayToItem[d]; ok {
			result = append(result, v)
		} else {
			result = append(result, model.UsageDaily{Date: d, UserIdentity: identity})
		}
	}
	return result, nil
}

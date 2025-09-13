package logic

import (
	"context"
	"fmt"
	"go-gin/internal/component/db"
	"go-gin/model"
	"time"
)

type AccountLogic struct{}

func NewAccountLogic() *AccountLogic { return &AccountLogic{} }

// UserPackageView 是前端所需的套餐视图，包含总量/已用/余额
type UserPackageView struct {
	PackageId      int64   `json:"package_id"`
	PackageName    string  `json:"package_name"`
	QuotaASRChars  int     `json:"quota_asr_chars"`
	QuotaTTSChars  int     `json:"quota_tts_chars"`
	RemainASRChars int     `json:"remain_asr_chars"`
	RemainTTSChars int     `json:"remain_tts_chars"`
	UsedASRChars   int     `json:"used_asr_chars"`
	UsedTTSChars   int     `json:"used_tts_chars"`
	ExpireAt       *string `json:"expire_at"`
}

func clampInt(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (l *AccountLogic) Profile(ctx context.Context, identity string) (map[string]any, error) {
	var au model.AccountUser
	_ = db.WithContext(ctx).Where("identity=?", identity).First(&au)
	// 附带返回用户专属声音ID（若无则为空字符串）
	voiceId := ""
	if identity != "" {
		if uv, err := model.NewUserVoiceModel().GetByMobile(ctx, identity); err == nil && uv != nil {
			voiceId = uv.VoiceId
		}
	}
	return map[string]any{
		"exists":       au.Id != 0,
		"identity":     identity,
		"status":       au.Status,
		"display_name": au.DisplayName,
		"voice_id":     voiceId,
	}, nil
}

func (l *AccountLogic) Packages(ctx context.Context, identity string) ([]UserPackageView, error) {
	// 读取用户的套餐余额记录
	var ups []model.UserPackage
	_ = db.WithContext(ctx).Where("user_identity=?", identity).Find(&ups)

	if len(ups) == 0 {
		return []UserPackageView{}, nil
	}

	// 收集 package_id 并查询对应套餐配额
	ids := make([]int64, 0, len(ups))
	for _, it := range ups {
		if it.PackageId != 0 {
			ids = append(ids, it.PackageId)
		}
	}

	pkgMap := map[int64]model.Package{}
	if len(ids) > 0 {
		var pkgs []model.Package
		_ = db.WithContext(ctx).Where("id IN ?", ids).Find(&pkgs)
		for _, p := range pkgs {
			pkgMap[p.Id] = p
		}
	}

	// 组装视图数据并做边界校验
	views := make([]UserPackageView, 0, len(ups))
	for _, up := range ups {
		p := pkgMap[up.PackageId]
		quotaASR := p.QuotaASRChars
		quotaTTS := p.QuotaTTSChars

		// 校验余额不能为负
		remainASR := up.RemainASRChars
		if remainASR < 0 {
			remainASR = 0
		}
		remainTTS := up.RemainTTSChars
		if remainTTS < 0 {
			remainTTS = 0
		}

		// 若总配额缺失（历史数据或未设置），以余额为上限作为总配额，避免进度条失真
		if quotaASR <= 0 {
			quotaASR = remainASR
		}
		if quotaTTS <= 0 {
			quotaTTS = remainTTS
		}

		// 进一步保证余额不超过总配额
		remainASR = clampInt(remainASR, 0, quotaASR)
		remainTTS = clampInt(remainTTS, 0, quotaTTS)

		usedASR := quotaASR - remainASR
		usedTTS := quotaTTS - remainTTS
		usedASR = clampInt(usedASR, 0, quotaASR)
		usedTTS = clampInt(usedTTS, 0, quotaTTS)

		views = append(views, UserPackageView{
			PackageId:      up.PackageId,
			PackageName:    p.Name,
			QuotaASRChars:  quotaASR,
			QuotaTTSChars:  quotaTTS,
			RemainASRChars: remainASR,
			RemainTTSChars: remainTTS,
			UsedASRChars:   usedASR,
			UsedTTSChars:   usedTTS,
			ExpireAt:       up.ExpireAt,
		})
	}

	return views, nil
}

// normalizeDate 将各种日期字符串规范为 yyyy-mm-dd
func normalizeDate(s string) string {
	if s == "" {
		return s
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.Format("2006-01-02")
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t.Format("2006-01-02")
	}
	// 已经是日期或前缀为日期
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func (l *AccountLogic) Usage(ctx context.Context, identity string) ([]model.UsageDaily, error) {
	var items []model.UsageDaily

	// 添加调试日志
	fmt.Printf("Usage query: identity=%s\n", identity)

	// 与写入口径保持一致：未登录/无身份时按 guest 统计
	if identity == "" {
		identity = "guest"
	}

	// 查询用户的使用记录
	err := db.WithContext(ctx).Where("user_identity=?", identity).Order("date desc").Limit(30).Find(&items).Error()
	if err != nil {
		fmt.Printf("Usage query error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Found %d usage records for identity=%s\n", len(items), identity)
	for _, item := range items {
		if item.ASRChars > 0 || item.TTSChars > 0 || item.Requests > 0 {
			fmt.Printf("Non-zero record: date=%s, asr=%d, tts=%d, requests=%d\n",
				item.Date, item.ASRChars, item.TTSChars, item.Requests)
		}
	}

	// 零填充近30天（按日期维度对齐）
	today := time.Now()
	dayToItem := map[string]model.UsageDaily{}
	for _, it := range items {
		key := normalizeDate(it.Date)
		dayToItem[key] = model.UsageDaily{
			Id:           it.Id,
			UserIdentity: it.UserIdentity,
			Date:         key,
			ASRChars:     it.ASRChars,
			TTSChars:     it.TTSChars,
			Requests:     it.Requests,
			CreatedAt:    it.CreatedAt,
		}
	}

	var result []model.UsageDaily
	for i := 0; i < 30; i++ {
		d := today.AddDate(0, 0, -i).Format("2006-01-02")
		if v, ok := dayToItem[d]; ok {
			result = append(result, v)
		} else {
			// 创建空记录
			result = append(result, model.UsageDaily{
				Date:         d,
				UserIdentity: identity,
				ASRChars:     0,
				TTSChars:     0,
				Requests:     0,
			})
		}
	}

	fmt.Printf("Returning %d total records (filled)\n", len(result))
	return result, nil
}

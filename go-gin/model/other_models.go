package model

type TTSHistory struct {
	Id           int64  `gorm:"column:id;primaryKey" json:"id"`
	UserIdentity string `gorm:"column:user_identity" json:"user_identity"`
	TextHash     string `gorm:"column:text_hash" json:"text_hash"`
	TextPreview  string `gorm:"column:text_preview" json:"text_preview"`
	CharCount    int    `gorm:"column:char_count" json:"char_count"`
	Speaker      string `gorm:"column:speaker" json:"speaker"`
	AudioUrl     string `gorm:"column:audio_url" json:"audio_url"`
	RequestId    string `gorm:"column:request_id" json:"request_id"`
	Status       int    `gorm:"column:status" json:"status"`
	CreatedAt    string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    string `gorm:"column:updated_at" json:"updated_at"`
}

func (TTSHistory) TableName() string { return "tts_history" }

type UserWhitelist struct {
	Id           int64  `gorm:"column:id;primaryKey" json:"id"`
	Identity     string `gorm:"column:identity" json:"identity"`
	IdentityType int    `gorm:"column:identity_type" json:"identity_type"`
	IsActive     int    `gorm:"column:is_active" json:"is_active"`
	Note         string `gorm:"column:note" json:"note"`
}

func (UserWhitelist) TableName() string { return "user_whitelist" }

type Package struct {
	Id            int64  `gorm:"column:id;primaryKey" json:"id"`
	Name          string `gorm:"column:name" json:"name"`
	QuotaASRChars int    `gorm:"column:quota_asr_chars" json:"quota_asr_chars"`
	QuotaTTSChars int    `gorm:"column:quota_tts_chars" json:"quota_tts_chars"`
	MonthlyReset  int    `gorm:"column:monthly_reset" json:"monthly_reset"`
}

func (Package) TableName() string { return "package" }

type UserPackage struct {
	Id             int64   `gorm:"column:id;primaryKey" json:"id"`
	UserIdentity   string  `gorm:"column:user_identity" json:"user_identity"`
	PackageId      int64   `gorm:"column:package_id" json:"package_id"`
	RemainASRChars int     `gorm:"column:remain_asr_chars" json:"remain_asr_chars"`
	RemainTTSChars int     `gorm:"column:remain_tts_chars" json:"remain_tts_chars"`
	ExpireAt       *string `gorm:"column:expire_at" json:"expire_at"`
}

func (UserPackage) TableName() string { return "user_package" }

type UsageDaily struct {
	Id           int64  `gorm:"column:id;primaryKey" json:"id"`
	UserIdentity string `gorm:"column:user_identity" json:"user_identity"`
	Date         string `gorm:"column:date" json:"date"`
	ASRChars     int    `gorm:"column:asr_chars" json:"asr_chars"`
	TTSChars     int    `gorm:"column:tts_chars" json:"tts_chars"`
	Requests     int    `gorm:"column:requests" json:"requests"`
}

func (UsageDaily) TableName() string { return "usage_daily" }

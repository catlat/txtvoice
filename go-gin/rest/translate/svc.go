package translate

import "context"

type ITranslateSvc interface {
	TranslateToZh(ctx context.Context, text string) (*TranslateResp, error)
}

type TranslateResp struct {
	Text      string `json:"text"`
	CharCount int    `json:"char_count"`
}

// Provider 枚举：当前支持 deepseek 与 bailian
type Provider string

const (
	ProviderDeepSeek Provider = "deepseek"
	ProviderBailian  Provider = "bailian"
)

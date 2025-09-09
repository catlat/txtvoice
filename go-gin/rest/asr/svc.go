package asr

import "context"

type IASRSvc interface {
	Recognize(ctx context.Context, audioUrl string) (*ASRResp, error)
}

// 火山引擎ASR响应数据结构
type ASRResp struct {
	Text      string `json:"text"`       // 合并后的完整文本
	CharCount int    `json:"char_count"` // 字符数统计
}

package translate

import (
	"go-gin/internal/httpc"
	"go-gin/util/jsonx"
)

// DeepSeek OpenAI兼容格式响应
type DeepSeekAPIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	// 错误信息（如果有）
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

var _ httpc.IResponse = (*DeepSeekAPIResponse)(nil)

func (r *DeepSeekAPIResponse) Parse(b []byte) error { return jsonx.Unmarshal(b, &r) }
func (r *DeepSeekAPIResponse) Valid() bool {
	return (len(r.Choices) > 0) || (r.Error != nil)
}
func (r *DeepSeekAPIResponse) IsSuccess() bool {
	return r.Error == nil && len(r.Choices) > 0 && r.Choices[0].Message.Content != ""
}
func (r *DeepSeekAPIResponse) Msg() string {
	if r.Error != nil {
		return r.Error.Message
	}
	return "success"
}
func (r *DeepSeekAPIResponse) ParseData() error {
	// 数据已直接解析到结构体中
	return nil
}


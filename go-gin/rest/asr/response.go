package asr

import (
	"go-gin/internal/httpc"
	"go-gin/util/jsonx"
)

// 火山引擎ASR API响应格式（根据实际响应示例）
type APIResponse struct {
	Header struct {
		ReqId   string `json:"reqid"`
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"header"`
	AudioInfo struct {
		Duration int64 `json:"duration"`
	} `json:"audio_info"`
	Result struct {
		Text       string `json:"text"`
		Utterances []struct {
			Text      string `json:"text"`
			StartTime int64  `json:"start_time"`
			EndTime   int64  `json:"end_time"`
			Words     []struct {
				Text       string `json:"text"`
				StartTime  int64  `json:"start_time"`
				EndTime    int64  `json:"end_time"`
				Confidence int    `json:"confidence"`
			} `json:"words"`
		} `json:"utterances"`
		Additions struct {
			Duration string `json:"duration"`
		} `json:"additions"`
	} `json:"result"`
	Data any `json:"data"` // 保留兼容性
}

var _ httpc.IResponse = (*APIResponse)(nil)

func (r *APIResponse) Parse(b []byte) error { return jsonx.Unmarshal(b, &r) }
func (r *APIResponse) Valid() bool          { return r.Header.Code != 0 || r.Result.Text != "" }        // 有状态码或有结果文本就是有效响应
func (r *APIResponse) IsSuccess() bool      { return r.Header.Code == 20000000 || r.Result.Text != "" } // 成功状态码或有结果文本
func (r *APIResponse) Msg() string          { return r.Header.Message }
func (r *APIResponse) ParseData() error {
	// 对于新格式，数据已经直接解析到Result字段中，无需额外处理
	return nil
}

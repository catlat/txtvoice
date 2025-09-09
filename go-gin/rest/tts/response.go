package tts

import (
	"go-gin/internal/httpc"
	"go-gin/util/jsonx"
)

type APIResponse struct {
	Code    *bool   `json:"success"`
	Message *string `json:"msg"`
	Data    any     `json:"param"`
}

var _ httpc.IResponse = (*APIResponse)(nil)

func (r *APIResponse) Parse(b []byte) error { return jsonx.Unmarshal(b, &r) }
func (r *APIResponse) Valid() bool          { return r.Code != nil && r.Message != nil }
func (r *APIResponse) IsSuccess() bool      { return r.Code != nil && *r.Code }
func (r *APIResponse) Msg() string {
	if r.Message == nil {
		return ""
	}
	return *r.Message
}
func (r *APIResponse) ParseData() error {
	dataStr, err := jsonx.Marshal(r.Data)
	if err != nil {
		return err
	}
	return jsonx.Unmarshal(dataStr, r.Data)
}


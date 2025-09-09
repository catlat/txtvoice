package dlyt

import (
	"encoding/json"
	"go-gin/internal/httpc"
	"go-gin/util/jsonx"
)

var (
	ApiResponseSuccessCode = true
)

// APIResponse 兼容 dlyt 服务返回格式，可同时支持：
// 1) 带包裹 {success,msg,param}
// 2) 直接返回数据对象（无 success/msg/param）
type APIResponse struct {
	Code    *bool           `json:"success"`
	Message *string         `json:"msg"`
	Param   json.RawMessage `json:"param"`
	Data    any             `json:"-"` // 目标数据指针，由调用方传入
}

var _ httpc.IResponse = (*APIResponse)(nil)

func (r *APIResponse) Parse(b []byte) error {
	// 先按包裹格式解析
	type envelope struct {
		Success *bool           `json:"success"`
		Msg     *string         `json:"msg"`
		Param   json.RawMessage `json:"param"`
	}
	var e envelope
	if err := jsonx.Unmarshal(b, &e); err == nil && (e.Success != nil || e.Msg != nil || len(e.Param) > 0) {
		r.Code = e.Success
		r.Message = e.Msg
		r.Param = e.Param
		return nil
	}
	// 兼容：无包裹，直接返回对象
	r.Param = json.RawMessage(b)
	ok := ApiResponseSuccessCode
	r.Code = &ok
	empty := ""
	r.Message = &empty
	return nil
}

func (r *APIResponse) Valid() bool {
	// 只要有 Code/Message 或有 Param 即视为有效
	return r.Code != nil || len(r.Param) > 0
}

func (r *APIResponse) IsSuccess() bool {
	if r.Code == nil {
		return true
	}
	return *r.Code == ApiResponseSuccessCode
}

func (r *APIResponse) Msg() string {
	if r.Message == nil {
		return ""
	}
	return *r.Message
}

func (r *APIResponse) ParseData() error {
	if r.Data == nil {
		return nil
	}
	// 将 param（或原始数据）解到 Data 指向的结构体
	return json.Unmarshal(r.Param, r.Data)
}

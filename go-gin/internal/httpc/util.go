package httpc

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
)

type M map[string]string

func NewClient() *Client {
	client := &Client{
		base: resty.New(),
	}
	client.SetTimeout(3 * time.Minute)

	// resty 默认启用连接复用，无需额外配置

	client.base.OnBeforeRequest(LogBeforeRequest)
	client.base.OnError(LogErrorHook)
	client.base.OnSuccess(LogSuccessHook)
	client.base.OnPanic(LogErrorHook)
	client.base.OnInvalid(LogErrorHook)
	return client
}

// NewStreamingClient 创建一个用于流式请求的客户端，不设置响应钩子以避免干扰流读取
func NewStreamingClient() *Client {
	client := &Client{
		base: resty.New(),
	}
	client.SetTimeout(3 * time.Minute)

	// resty 默认启用连接复用，匹配服务端keep-alive时间（1分钟）

	// 只保留请求前的日志，避免响应钩子干扰流式读取
	client.base.OnBeforeRequest(LogBeforeRequest)
	return client
}

func GET(ctx context.Context, url string) *Request {
	return NewClient().
		NewRequest().
		GET(url).
		SetContext(ctx)
}

func POST(ctx context.Context, url string) *Request {
	return NewClient().
		NewRequest().
		POST(url).
		SetContext(ctx)
}

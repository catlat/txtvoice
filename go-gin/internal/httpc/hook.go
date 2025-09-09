package httpc

import (
	"context"
	"go-gin/internal/component/logx"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

func LogBeforeRequest(c *resty.Client, r *resty.Request) error {
	url := r.URL
	if !strings.HasPrefix(url, "http") {
		url = c.BaseURL + url
	}
	logx.RestyLoggerInstance.Info().Ctx(r.Context()).
		Str("keywords", "request").
		Str("url", url).
		Str("method", r.Method).
		Any("header", r.Header).
		Str("query", r.QueryParam.Encode()).
		Any("post", r.FormData.Encode()).
		Any("body", r.Body).
		Send()
	return nil
}

func LogErrorHook(r *resty.Request, err error) {
	if responseErr, ok := err.(*resty.ResponseError); ok {
		LogResponse(r.Context(), responseErr.Response)
	}
	logx.RestyLoggerInstance.Info().Ctx(r.Context()).
		Str("keywords", "error hook").
		Str("msg", err.Error()).
		Send()
}

func LogSuccessHook(c *resty.Client, r *resty.Response) {
	LogResponse(r.Request.Context(), r)
}

func LogResponse(ctx context.Context, r *resty.Response) {
	e := logx.RestyLoggerInstance.
		Info().
		Ctx(ctx).
		Str("keywords", "response")

	// 避免对流式响应执行 r.String() 导致底层 RawBody 被读尽或关闭
	url := r.Request.URL
	logBody := true
	if strings.Contains(url, "/api/v3/tts/unidirectional") || r.Request.Header.Get("X-Log-Body") == "ignore" {
		logBody = false
	}
	if logBody {
		e = e.Str("body", r.String())
	} else {
		e = e.Str("body", "<streaming body omitted>")
	}

	if r.StatusCode() != http.StatusOK {
		e = e.Int("status", r.StatusCode())
	}
	if r.Error() != nil {
		e = e.Any("error", r.Error())
	}
	e.Send()
}

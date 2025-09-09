package middleware

import "go-gin/internal/httpx"

func Init(r *httpx.Engine) {
	// CORS 应该最先处理
	r.Before(CORS())

	r.Before(BeforeSampleA(), BeforeSampleB())
	r.After(AfterSampleB())
	// r.Before(TokenCheck())

}

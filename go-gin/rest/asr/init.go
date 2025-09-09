package asr

var (
	Svc IASRSvc = (*ASRSvc)(nil)
)

func Init(url string) {
	Svc = NewASRSvc(url)
}


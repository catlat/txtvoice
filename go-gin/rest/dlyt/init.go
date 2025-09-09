package dlyt

var (
	Svc IYtSvc = (*YtSvc)(nil)
)

func Init(url string) {
	if url == "" || url == "local" {
		Svc = NewLocalYtSvc()
		return
	}
	Svc = NewYtSvc(url)
}

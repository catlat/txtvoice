package translate

var (
	Svc ITranslateSvc = (*TranslateSvc)(nil)
)

func Init(url string) { Svc = NewTranslateSvc(url) }


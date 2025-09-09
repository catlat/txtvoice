package tts

var (
	Svc ITTSSvc = (*TTSSvc)(nil)
)

func Init(url string) { Svc = NewTTSSvc(url) }


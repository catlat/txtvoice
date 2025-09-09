package asr

type VolcCreds struct {
	AppId         string
	AccessKey     string
	ASRResourceId string
}

var volcCreds VolcCreds

func SetVolcCreds(c VolcCreds) { volcCreds = c }

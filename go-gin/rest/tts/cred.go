package tts

type VolcCreds struct {
	AppId         string
	AccessKey     string
	TTSResourceId string
}

var volcCreds VolcCreds

func SetVolcCreds(c VolcCreds) { volcCreds = c }

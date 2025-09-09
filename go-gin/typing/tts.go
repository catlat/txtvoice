package typing

type TTSSynthesizeReq struct {
	Text    string `form:"text" binding:"required" label:"文本"`
	Speaker string `form:"speaker" binding:"required" label:"说话人"`
}

type TTSSynthesizeReply struct {
	AudioUrl string `json:"audio_url"`
}


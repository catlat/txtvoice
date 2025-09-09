package typing

type YtInfoReq struct {
	IdOrUrl  string `form:"id_or_url" json:"id_or_url" binding:"required" label:"视频ID或链接"`
	Platform string `form:"platform" json:"platform" binding:"omitempty" label:"平台类型"`
}

type YtInfoReply struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	DurationSec  int    `json:"duration_sec"`
	Views        int64  `json:"views"`
	PublishDate  string `json:"publish_date"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type YtTextReq struct {
	IdOrUrl   string `form:"id_or_url" json:"id_or_url" binding:"required" label:"视频ID或链接"`
	TargetLan string `form:"target_lang" json:"target_lang" binding:"omitempty" label:"目标语言"`
	Platform  string `form:"platform" json:"platform" binding:"omitempty" label:"平台类型"`
}

type YtTextReply struct {
	OriginalText   string `json:"original_text"`
	TranslatedText string `json:"translated_text"`
	AudioData      string `json:"audio_data"`
	AudioType      string `json:"audio_type"` // 音频文件类型，如：mp3, m4a, wav, aac, ogg, flac
}

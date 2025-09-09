package model

type YoutubeVideo struct {
	Id           int64   `gorm:"column:id;primaryKey" json:"id"`
	SourceSite   string  `gorm:"column:source_site" json:"source_site"`
	VideoId      string  `gorm:"column:video_id" json:"video_id"`
	Title        string  `gorm:"column:title" json:"title"`
	Description  string  `gorm:"column:description" json:"description"`
	ChannelTitle string  `gorm:"column:channel_title" json:"channel_title"`
	DurationSec  int     `gorm:"column:duration_sec" json:"duration_sec"`
	PublishedAt  *string `gorm:"column:published_at" json:"published_at"`
	ThumbnailUrl string  `gorm:"column:thumbnail_url" json:"thumbnail_url"`
	AudioUrl     string  `gorm:"column:audio_url" json:"audio_url"`
	Status       int     `gorm:"column:status" json:"status"`
}

func (YoutubeVideo) TableName() string { return "youtube_video" }

type YoutubeTranscript struct {
	Id                 int64  `gorm:"column:id;primaryKey" json:"id"`
	VideoId            int64  `gorm:"column:video_id" json:"video_id"`
	Language           string `gorm:"column:language" json:"language"`
	OriginalText       string `gorm:"column:original_text" json:"original_text"`
	TranslatedText     string `gorm:"column:translated_text" json:"translated_text"`
	AsrCharCount       int    `gorm:"column:asr_char_count" json:"asr_char_count"`
	TranslateCharCount int    `gorm:"column:translate_char_count" json:"translate_char_count"`
}

func (YoutubeTranscript) TableName() string { return "youtube_transcript" }

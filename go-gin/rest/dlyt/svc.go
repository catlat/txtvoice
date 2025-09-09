package dlyt

import (
	"context"
)

type IYtSvc interface {
	Info(ctx context.Context, idOrUrl string) (*InfoResp, error)
	Audio(ctx context.Context, idOrUrl string) (*AudioResp, error)
}

type InfoResp struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	DurationSec  int    `json:"duration_sec"`
	Views        int64  `json:"views"`
	PublishDate  string `json:"publish_date"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type AudioResp struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	AudioUrl string `json:"audio_url"`
}

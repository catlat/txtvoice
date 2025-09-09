package dlyt

import (
	"context"
	"go-gin/const/errcode"
	"go-gin/internal/httpc"
)

const (
	InfoURL  = "/api/yt/info"
	AudioURL = "/api/yt/audio"
)

type YtSvc struct {
	httpc.BaseSvc
}

func NewYtSvc(url string) IYtSvc {
	return &YtSvc{BaseSvc: *httpc.NewBaseSvc(url)}
}

func (s *YtSvc) Info(ctx context.Context, idOrUrl string) (resp *InfoResp, err error) {
	result := APIResponse{Data: &resp}
	err = s.Client().
		NewRequest().
		SetContext(ctx).
		SetHeaders(map[string]string{"Content-Type": "application/json", "Accept": "application/json"}).
		POST(InfoURL).
		SetBody(map[string]string{"id_or_url": idOrUrl}).
		SetResult(&result).
		Exec()
	if err != nil {
		return nil, errcode.ErrDLYTUpstream
	}
	return resp, nil
}

func (s *YtSvc) Audio(ctx context.Context, idOrUrl string) (resp *AudioResp, err error) {
	result := APIResponse{Data: &resp}
	err = s.Client().
		NewRequest().
		SetContext(ctx).
		SetHeaders(map[string]string{"Content-Type": "application/json", "Accept": "application/json"}).
		POST(AudioURL).
		SetBody(map[string]string{"id_or_url": idOrUrl}).
		SetResult(&result).
		Exec()
	if err != nil {
		return nil, errcode.ErrDLYTUpstream
	}
	return resp, nil
}

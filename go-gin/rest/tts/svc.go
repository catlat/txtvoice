package tts

import "context"

type ITTSSvc interface {
	Synthesize(ctx context.Context, text, speaker string) (*TTSResp, error)
	// SynthesizeWithResource performs TTS using the specified resource id with no fallback logic
	SynthesizeWithResource(ctx context.Context, text, speaker, resourceId string) (*TTSResp, error)
}

type TTSResp struct {
	Audio          []byte `json:"-"`
	AudioUrl       string `json:"audio_url"`
	Size           int    `json:"-"`
	UsedSpeaker    string `json:"-"`
	UsedResourceId string `json:"-"`
	RequestId      string `json:"-"`
}

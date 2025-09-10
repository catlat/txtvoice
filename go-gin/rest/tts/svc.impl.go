package tts

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"go-gin/const/errcode"
	"go-gin/internal/httpc"
	"io"
	"log"
	"strings"

	"github.com/google/uuid"
)

const (
	SynthesizeURL   = "/api/v3/tts/unidirectional"
	DefaultSpeaker  = "zh_female_shuangkuaisisi_moon_bigtts"
	DefaultResource = "volc.service_type.10029"
)

type TTSSvc struct {
	httpc.BaseSvc
	baseURL string
}

func NewTTSSvc(url string) ITTSSvc {
	// 如果传的是完整HTTP地址，避免设置到Client BaseURL，防止拼接重复
	if strings.HasPrefix(url, "http") {
		return &TTSSvc{BaseSvc: *httpc.NewStreamingBaseSvc(""), baseURL: url}
	}
	return &TTSSvc{BaseSvc: *httpc.NewStreamingBaseSvc(url), baseURL: url}
}

type ttsStreamResp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Data     string `json:"data"`
	Sentence *struct {
		Phonemes []interface{} `json:"phonemes"`
		Text     string        `json:"text"`
		Words    []struct {
			Confidence float64 `json:"confidence"`
			EndTime    float64 `json:"endTime"`
			StartTime  float64 `json:"startTime"`
			Word       string  `json:"word"`
		} `json:"words"`
	} `json:"sentence,omitempty"`
}

func (s *TTSSvc) Synthesize(ctx context.Context, text, speaker string) (resp *TTSResp, err error) {
	// 第一次按用户入参尝试；失败(资源不匹配)则降级到默认普通音色+资源
	if r, e := s.doOnce(ctx, text, speaker, pickResourceBySpeaker(speaker)); e == nil {
		return r, nil
	}
	// fallback
	return s.doOnce(ctx, text, DefaultSpeaker, defaultResourceId())
}

// SynthesizeWithResource 使用明确的资源ID进行合成，不做任何回退
func (s *TTSSvc) SynthesizeWithResource(ctx context.Context, text, speaker, resourceId string) (*TTSResp, error) {
	return s.doOnce(ctx, text, speaker, resourceId)
}

func pickResourceBySpeaker(speaker string) string {
	sp := strings.ToLower(speaker)
	if strings.HasPrefix(sp, "rec_") || sp == "custom_mix_bigtts" {
		return "volc.megatts.default"
	}
	return defaultResourceId()
}

// defaultResourceId 返回配置的 TTSResourceId；若未配置则回退到编译期默认值
func defaultResourceId() string {
	if volcCreds.TTSResourceId != "" {
		return volcCreds.TTSResourceId
	}
	return DefaultResource
}

// detectExplicitLanguage 根据音色推断主要语种，而非文本语种
func detectExplicitLanguage(speaker string) string {
	sp := strings.ToLower(speaker)
	if strings.HasPrefix(sp, "zh_") || strings.Contains(sp, "chinese") {
		return "zh"
	} else if strings.HasPrefix(sp, "en_") || strings.Contains(sp, "english") {
		return "en"
	} else if strings.HasPrefix(sp, "ja_") || strings.Contains(sp, "japanese") {
		return "ja"
	}
	// 复刻音色或未知音色，启用多语种前端
	return "zh,en,ja,es-mx,id,pt-br,de,fr"
}

func (s *TTSSvc) doOnce(ctx context.Context, text, speaker, resourceId string) (*TTSResp, error) {
	// 验证凭据
	if volcCreds.AppId == "" || volcCreds.AccessKey == "" {
		log.Printf("TTS ERROR: Missing Volc credentials")
		return nil, errcode.ErrTTSUpstream
	}

	// 根据音色设置语种，启用语言检测以支持跨语种合成
	additions := map[string]any{
		"enable_language_detector": true,
		"explicit_language":        detectExplicitLanguage(speaker),
		"disable_markdown_filter":  true,
	}

	// additions 必须是 JSON 字符串，不是对象
	additionsJSON, _ := json.Marshal(additions)

	payload := map[string]any{
		"user": map[string]any{"uid": volcCreds.AppId},
		"req_params": map[string]any{
			"text":    text,
			"speaker": speaker,
			"audio_params": map[string]any{
				"format":           "mp3",
				"sample_rate":      24000,
				"enable_timestamp": true,
			},
			"additions": string(additionsJSON),
		},
	}
	endpoint := SynthesizeURL
	if strings.HasPrefix(s.baseURL, "http") {
		endpoint = s.baseURL
	}
	reqID := uuid.New().String()
	req := s.Client().NewRequest().SetContext(ctx).POST(endpoint).
		SetHeaders(map[string]string{
			"X-Api-App-Id":      volcCreds.AppId,
			"X-Api-Access-Key":  volcCreds.AccessKey,
			"X-Api-Resource-Id": resourceId,
			"X-Api-Request-Id":  reqID,
			"Content-Type":      "application/json",
			"Connection":        "keep-alive",
		}).
		SetBody(payload).
		SetDoNotParseResponse(true)

	res, e := req.Send()
	if e != nil {
		log.Printf("TTS request failed: %v (speaker=%s, resource=%s)", e, speaker, resourceId)
		return nil, errcode.ErrTTSUpstream
	}
	defer func() { _ = res.RawBody().Close() }()
	reader := bufio.NewReader(res.RawBody())
	var audio []byte
	var sawMismatch bool
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("TTS Synthesize error 2:", err)
			return nil, errcode.ErrTTSUpstream
		}
		if len(line) == 0 {
			continue
		}
		var frame ttsStreamResp
		if json.Unmarshal(line, &frame) != nil {
			continue
		}
		if frame.Data != "" {
			chunk, _ := base64.StdEncoding.DecodeString(frame.Data)
			audio = append(audio, chunk...)
		}
		if frame.Sentence != nil && frame.Sentence.Text != "" {
			log.Printf("TTS sentence: %s (words: %d)", frame.Sentence.Text, len(frame.Sentence.Words))
		}
		if frame.Code == 20000000 {
			break
		}
		if frame.Code == 55000000 && strings.Contains(strings.ToLower(frame.Message), "resource") {
			sawMismatch = true
			break
		}
	}
	if sawMismatch {
		log.Printf("TTS resource mismatch: speaker=%s, resource=%s", speaker, resourceId)
		return nil, errcode.ErrTTSUpstream
	}
	// 若未收到任何音频分片，视为失败，让上层触发回退逻辑
	if len(audio) == 0 {
		log.Printf("TTS empty audio: speaker=%s, resource=%s", speaker, resourceId)
		return nil, errcode.ErrTTSUpstream
	}
	log.Printf("TTS success: speaker=%s, resource=%s, bytes=%d", speaker, resourceId, len(audio))
	return &TTSResp{Audio: audio, Size: len(audio), UsedSpeaker: speaker, UsedResourceId: resourceId, RequestId: reqID}, nil
}

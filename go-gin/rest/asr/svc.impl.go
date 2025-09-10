package asr

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go-gin/const/errcode"
	"go-gin/internal/httpc"
	"go-gin/internal/traceid"
)

const (
	FlashURL = "/api/v3/auc/bigmodel/recognize/flash"
)

type ASRSvc struct{ httpc.BaseSvc }

func NewASRSvc(url string) IASRSvc { return &ASRSvc{BaseSvc: *httpc.NewBaseSvc(url)} }

func (s *ASRSvc) Recognize(ctx context.Context, audioUrl string) (resp *ASRResp, err error) {
	log.Printf("[ASR] 开始语音识别, audioUrl: %s", audioUrl)

	// 检查认证信息
	if volcCreds.AppId == "" || volcCreds.AccessKey == "" || volcCreds.ASRResourceId == "" {
		log.Printf("[ASR] 认证信息不完整 - AppId: %s, AccessKey: %s, ASRResourceId: %s",
			maskString(volcCreds.AppId), maskString(volcCreds.AccessKey), maskString(volcCreds.ASRResourceId))
		return nil, errcode.ErrASRUpstream
	}

	requestId := traceid.New()
	var result APIResponse

	// 构建请求体
	var body map[string]any
	if isLocalStatic(audioUrl) {
		// 本地文件模式：读取文件并转换为Base64
		localPath := mapStaticToLocal(audioUrl)
		log.Printf("[ASR] 读取本地文件转Base64 - path: %s", localPath)

		fileData, readErr := os.ReadFile(localPath)
		if readErr != nil {
			log.Printf("[ASR] 读取本地文件失败 - path: %s, error: %v", localPath, readErr)
			return nil, errcode.ErrASRUpstream
		}

		base64Data := base64.StdEncoding.EncodeToString(fileData)
		log.Printf("[ASR] 文件转Base64成功 - size: %d bytes, base64Length: %d", len(fileData), len(base64Data))

		body = map[string]any{
			"user": map[string]any{
				"uid": volcCreds.AppId,
			},
			"audio": map[string]any{
				"data": base64Data,
			},
			"request": map[string]any{
				"model_name": "bigmodel",
				"language":   "en-US",
				"enable_ddc": true,
			},
		}
	} else {
		// URL模式
		body = map[string]any{
			"user": map[string]any{
				"uid": volcCreds.AppId,
			},
			"audio": map[string]any{
				"url": audioUrl,
			},
			"request": map[string]any{
				"model_name": "bigmodel",
				"language":   "en-US",
				"enable_ddc": true,
			},
		}
	}

	log.Printf("[ASR] 发送请求 - RequestId: %s, AudioSource: %s", requestId,
		func() string {
			if isLocalStatic(audioUrl) {
				return "base64_file"
			} else {
				return "url"
			}
		}())

	// 发送请求
	if execErr := s.Client().NewRequest().SetContext(ctx).POST(FlashURL).
		SetHeaders(map[string]string{
			"X-Api-App-Key":     volcCreds.AppId,
			"X-Api-Access-Key":  volcCreds.AccessKey,
			"X-Api-Resource-Id": volcCreds.ASRResourceId,
			"X-Api-Sequence":    "-1",
			"X-Api-Request-Id":  requestId,
		}).
		SetBody(body).
		SetResult(&result).
		Exec(); execErr != nil {
		log.Printf("[ASR] 请求失败 - RequestId: %s, Error: %v", requestId, execErr)
		return nil, errcode.ErrASRUpstream
	}

	// 检查API响应
	if !result.Valid() {
		log.Printf("[ASR] 响应格式无效 - RequestId: %s, Response: %+v", requestId, result)
		return nil, errcode.ErrASRUpstream
	}
	if !result.IsSuccess() {
		log.Printf("[ASR] API返回失败 - RequestId: %s, Code: %d, Message: %s",
			requestId, result.Header.Code, result.Msg())
		return nil, errcode.ErrASRUpstream
	}

	// 转换火山引擎数据为统一格式
	resp = &ASRResp{}

	// 优先使用result.text，如果为空则从utterances合并
	if result.Result.Text != "" {
		resp.Text = strings.TrimSpace(result.Result.Text)
	} else {
		var fullText strings.Builder
		for _, utterance := range result.Result.Utterances {
			fullText.WriteString(utterance.Text)
		}
		resp.Text = strings.TrimSpace(fullText.String())
	}

	resp.CharCount = len([]rune(resp.Text)) // 使用rune计算字符数（支持中文）

	if resp.Text == "" {
		log.Printf("[ASR] 识别结果为空 - RequestId: %s", requestId)
		return nil, errcode.ErrASRUpstream
	}

	log.Printf("[ASR] 识别成功 - RequestId: %s, CharCount: %d, TextLength: %d, Utterances: %d, Duration: %dms",
		requestId, resp.CharCount, len(resp.Text), len(result.Result.Utterances), result.AudioInfo.Duration)
	log.Printf("[ASR] 识别文本预览: %.200s...", resp.Text)

	return resp, nil
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}

func isLocalStatic(u string) bool {
	u = strings.TrimSpace(u)
	return strings.HasPrefix(u, "/static/")
}

func mapStaticToLocal(u string) string {
	// /static/yt/audio/xxx -> ./public/yt/audio/xxx
	trimmed := strings.TrimPrefix(u, "/static/")
	return filepath.Join("public", filepath.FromSlash(trimmed))
}

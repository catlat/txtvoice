package translate

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go-gin/const/errcode"
	"go-gin/internal/httpc"
	"go-gin/internal/traceid"
)

const (
	ChatCompletionsURL        = "/chat/completions"
	BailianChatCompletionsURL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
)

type TranslateSvc struct {
	httpc.BaseSvc
}

var deepseekCreds DeepSeekCreds
var bailianCreds BailianCreds
var provider Provider = ProviderDeepSeek

type DeepSeekCreds struct {
	ApiKey string
}

type BailianCreds struct {
	ApiKey string
}

// SetDeepSeekCreds 设置DeepSeek凭据
func SetDeepSeekCreds(creds DeepSeekCreds) {
	deepseekCreds = creds
	log.Printf("[Translate] DeepSeek凭据已设置 - ApiKey: %s", maskApiKey(creds.ApiKey))
}

// SetBailianCreds 设置阿里云百炼凭据
func SetBailianCreds(creds BailianCreds) {
	bailianCreds = creds
	log.Printf("[Translate] Bailian凭据已设置 - ApiKey: %s", maskApiKey(creds.ApiKey))
}

// SetProvider 设置翻译提供方
func SetProvider(p Provider) { provider = p }

func NewTranslateSvc(url string) ITranslateSvc {
	// 直接使用DeepSeek API地址
	return &TranslateSvc{BaseSvc: *httpc.NewBaseSvc("https://api.deepseek.com")}
}

func (s *TranslateSvc) TranslateToZh(ctx context.Context, text string) (resp *TranslateResp, err error) {
	log.Printf("[Translate] 开始翻译 - 文本长度: %d", len(text))

	// 检查API Key（根据 provider）
	switch provider {
	case ProviderBailian:
		if bailianCreds.ApiKey == "" {
			log.Printf("[Translate] Bailian API Key未配置")
			return nil, errcode.ErrTranslateUp
		}
	default:
		if deepseekCreds.ApiKey == "" {
			log.Printf("[Translate] DeepSeek API Key未配置")
			return nil, errcode.ErrTranslateUp
		}
	}

	requestId := traceid.New()

	// 构建OpenAI兼容的请求体
	model := "deepseek-chat"
	if provider == ProviderBailian {
		model = "qwen-turbo-latest"
	}
	requestBody := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "你是一个专业的翻译助手。请将用户提供的英文文本翻译成中文，保持原意和语言风格。请直接返回翻译后的中文文本，不要添加任何解释或前缀。",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("请翻译以下英文文本为中文：\n\n%s", text),
			},
		},
		"stream": false,
	}

	var result DeepSeekAPIResponse
	req := s.Client().NewRequest().SetContext(ctx)
	if provider == ProviderBailian {
		// 直连百炼完整URL（不依赖 BaseURL）
		req = req.POST(BailianChatCompletionsURL).SetHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", bailianCreds.ApiKey),
			"Content-Type":  "application/json",
		})
	} else {
		// 使用 DeepSeek BaseURL + 相对路径
		req = req.POST(ChatCompletionsURL).SetHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", deepseekCreds.ApiKey),
			"Content-Type":  "application/json",
		})
	}
	err = req.SetBody(requestBody).SetResult(&result).Exec()

	if err != nil {
		log.Printf("[Translate] 请求失败 - RequestId: %s, Error: %v", requestId, err)
		return nil, errcode.ErrTranslateUp
	}

	// 检查响应
	if !result.Valid() {
		log.Printf("[Translate] 响应格式无效 - RequestId: %s, Response: %+v", requestId, result)
		return nil, errcode.ErrTranslateUp
	}

	if !result.IsSuccess() {
		log.Printf("[Translate] API返回失败 - RequestId: %s, Error: %s", requestId, result.Msg())
		return nil, errcode.ErrTranslateUp
	}

	// 提取翻译结果
	translatedText := ""
	if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
		content := strings.TrimSpace(result.Choices[0].Message.Content)

		// 移除常见的前后缀和格式标记
		content = strings.TrimPrefix(content, "翻译结果：")
		content = strings.TrimPrefix(content, "翻译为：")
		content = strings.TrimPrefix(content, "中文翻译：")
		content = strings.TrimPrefix(content, "以下是翻译：")

		// 移除引号包围
		if strings.HasPrefix(content, `"`) && strings.HasSuffix(content, `"`) {
			content = content[1 : len(content)-1]
		}

		translatedText = strings.TrimSpace(content)

		log.Printf("[Translate] 提取翻译文本成功 - RequestId: %s, 长度: %d", requestId, len(translatedText))
	}

	if translatedText == "" {
		log.Printf("[Translate] 提取翻译文本失败 - RequestId: %s, Content: %s", requestId, result.Choices[0].Message.Content)
		return nil, errcode.ErrTranslateUp
	}

	resp = &TranslateResp{
		Text:      translatedText,
		CharCount: len([]rune(translatedText)), // 使用rune计算字符数
	}

	log.Printf("[Translate] 翻译成功 - RequestId: %s, CharCount: %d, UsageTokens: %d",
		requestId, resp.CharCount, result.Usage.TotalTokens)
	log.Printf("[Translate] 翻译文本预览: %.200s...", resp.Text)

	return resp, nil
}

func maskApiKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}

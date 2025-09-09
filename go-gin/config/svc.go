package config

import (
	"log"

	"go-gin/internal/qiniu"
	"go-gin/rest/asr"
	"go-gin/rest/dlyt"
	"go-gin/rest/login"
	"go-gin/rest/mylogin"
	"go-gin/rest/translate"
	"go-gin/rest/tts"
	"go-gin/rest/user"
)

type SvcConfig struct {
	UserSvcUrl        string `yaml:"user_url"`
	LoginSvcUrl       string `yaml:"login_url"`
	ASRUrl            string `yaml:"asr_url"`
	TTSUrl            string `yaml:"tts_url"`
	TranslateProvider string `yaml:"translate_provider"`
}

func InitSvc() {
	svcConfig := instance.Svc
	// 初始化第三方请求服务
	user.Init(svcConfig.UserSvcUrl)
	login.Init(svcConfig.LoginSvcUrl)
	mylogin.Init(svcConfig.LoginSvcUrl)
	// dlyt 始终使用本地实现
	dlyt.Init("")
	asr.Init(svcConfig.ASRUrl)
	translate.Init("") // URL在service内部写死
	tts.Init(svcConfig.TTSUrl)

	// 注入火山凭据
	volc := instance.Creds.Volc
	asr.SetVolcCreds(asr.VolcCreds{AppId: volc.AppId, AccessKey: volc.AccessKey, ASRResourceId: volc.ASRResourceId})
	tts.SetVolcCreds(tts.VolcCreds{AppId: volc.AppId, AccessKey: volc.AccessKey, TTSResourceId: volc.TTSResourceId})

	// 注入DeepSeek凭据
	ds := instance.Creds.Deepseek
	translate.SetDeepSeekCreds(translate.DeepSeekCreds{ApiKey: ds.ApiKey})
	log.Printf("deepseek config: api_key=%s", maskKey(ds.ApiKey))

	// 注入百炼凭据
	bl := instance.Creds.Bailian
	translate.SetBailianCreds(translate.BailianCreds{ApiKey: bl.ApiKey})
	log.Printf("bailian config: api_key=%s", maskKey(bl.ApiKey))

	// 选择翻译提供方（默认 deepseek），可在 yaml: svc.translate_provider=bailian 切换
	provider := translate.ProviderDeepSeek
	if svcConfig.TranslateProvider == "bailian" {
		provider = translate.ProviderBailian
	}
	translate.SetProvider(provider)

	// 注入七牛凭据（保持与 dlyt 一致的参数命名）
	q := instance.Creds.Qiniu
	log.Printf("qiniu config: access_key=%s secret_key=%s bucket=%s domain=%s",
		maskKey(q.AccessKey), maskKey(q.SecretKey), q.Bucket, q.Domain)
	if q.AccessKey != "" && q.SecretKey != "" && q.Bucket != "" {
		dlyt.SetQiniuConfig(qiniu.Config{AccessKey: q.AccessKey, SecretKey: q.SecretKey, Bucket: q.Bucket, Domain: q.Domain})
		log.Printf("qiniu uploader initialized successfully")
		log.Printf("qiniu domain configured: %s", q.Domain)
	} else {
		log.Printf("qiniu config incomplete, uploader disabled - missing: access_key=%v secret_key=%v bucket=%v",
			q.AccessKey == "", q.SecretKey == "", q.Bucket == "")
		log.Printf("WARNING: 音频无法上传到七牛，ASR识别可能失败")
	}
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}

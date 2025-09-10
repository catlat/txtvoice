package dlyt

// Options 为 dlyt 本地服务的运行选项，由上层注入，避免与 config 产生循环依赖
type Options struct {
	// Bilibili 音频处理模式：local | url（默认 local）
	BilibiliAudioMode string
	// Bilibili URL 策略：raw | proxy（当前仅 raw 生效，占位）
	BilibiliURLStrategy string
}

// 默认：B站直链（url 模式），策略 raw（占位）
var pkgOptions = Options{BilibiliAudioMode: "url", BilibiliURLStrategy: "raw"}

// SetOptions 由上层（config.InitSvc）在启动时调用
func SetOptions(opt Options) {
	// 仅当有值时覆盖默认值
	if opt.BilibiliAudioMode != "" {
		pkgOptions.BilibiliAudioMode = opt.BilibiliAudioMode
	}
	if opt.BilibiliURLStrategy != "" {
		pkgOptions.BilibiliURLStrategy = opt.BilibiliURLStrategy
	}
}

var (
	Svc IYtSvc = (*YtSvc)(nil)
)

func Init(url string) {
	if url == "" || url == "local" {
		Svc = NewLocalYtSvc()
		return
	}
	Svc = NewYtSvc(url)
}

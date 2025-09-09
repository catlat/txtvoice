package config

type CredsConfig struct {
	Deepseek struct {
		ApiKey string `yaml:"api_key"`
	} `yaml:"deepseek"`

	Bailian struct {
		ApiKey string `yaml:"api_key"`
	} `yaml:"bailian"`

	Volc struct {
		AppId         string `yaml:"app_id"`
		AccessKey     string `yaml:"access_key"`
		ASRResourceId string `yaml:"asr_resource_id"`
		TTSResourceId string `yaml:"tts_resource_id"`
	} `yaml:"volc"`

	Qiniu struct {
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Bucket    string `yaml:"bucket"`
		Domain    string `yaml:"domain"`
	} `yaml:"qiniu"`
}

func GetCreds() CredsConfig { return instance.Creds }

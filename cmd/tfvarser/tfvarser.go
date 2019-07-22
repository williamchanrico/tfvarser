package tfvarser

import "os"

// Flags contains run-time flags
type Flags struct {
	Provider    string
	ProviderObj string
	LimitNames  string
	LimitIDs    string
	TemplateDir string
}

// Config contains tfvars config
type Config struct {
	AlicloudAccessKey string `envconfig:"ALICLOUD_ACCESS_KEY" required:"true"`
	AlicloudSecretKey string `envconfig:"ALICLOUD_SECRET_KEY" required:"true"`
	AlicloudRegionID  string `default:"ap-southeast-1"`
}

// Run is the entrypoint for main autoscaleapp process
func Run(appFlags *Flags, cfg Config) (int, error) {
	if _, err := os.Stat(appFlags.TemplateDir); os.IsNotExist(err) {
		return 1, err
	}

	switch appFlags.Provider {
	case "ali":
		return aliyunProvider(appFlags, cfg)
	default:
		return 1, ErrProviderNotSupported
	}
}

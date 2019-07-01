package aliyun

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/williamchanrico/tfvarser/aliyun/ecs"
	"github.com/williamchanrico/tfvarser/aliyun/ess"
)

// Config contains aliyun client config
type Config struct {
	AccessKey string
	SecretKey string
	RegionID  string
}

// Client is aliyun api struct
type Client struct {
	ESS *ess.Client
	ECS *ecs.Client
}

// New returns new aliyun api client
func New(cfg Config) (*Client, error) {
	creds := credentials.NewAccessKeyCredential(cfg.AccessKey, cfg.SecretKey)
	config := sdk.NewConfig().
		WithAutoRetry(true).
		WithMaxRetryTime(3).
		WithTimeout(10 * time.Second).
		WithEnableAsync(true)
	acsClient, err := sdk.NewClientWithOptions(cfg.RegionID, config, creds)
	acsClient.EnableAsync(5, 100)
	if err != nil {
		return nil, err
	}

	return &Client{
		ESS: ess.New(*acsClient),
		ECS: ecs.New(*acsClient),
	}, nil
}

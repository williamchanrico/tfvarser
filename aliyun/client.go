package aliyun

import (
	ecssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
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
	ecsClient, err := ecssdk.NewClientWithAccessKey(
		cfg.RegionID,
		cfg.AccessKey,
		cfg.SecretKey,
	)
	ecsClient.EnableAsync(5, 10)
	if err != nil {
		return nil, err
	}
	ecs := ecs.New(ecsClient)

	essClient, err := esssdk.NewClientWithAccessKey(
		cfg.RegionID,
		cfg.AccessKey,
		cfg.SecretKey,
	)
	essClient.EnableAsync(5, 10)
	if err != nil {
		return nil, err
	}
	ess := ess.New(essClient)

	return &Client{
		ESS: ess,
		ECS: ecs,
	}, nil
}

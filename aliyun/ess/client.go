package ess

import (
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// Client is the ess client
type Client struct {
	ess *esssdk.Client
}

// Config contains ess client config
type Config struct {
	AccessKey string
	SecretKey string
	RegionID  string
}

// New returns a new ess client
func New(c *Config) (*Client, error) {
	// Create an ESS client
	essClient, err := esssdk.NewClientWithAccessKey(
		c.RegionID,
		c.AccessKey,
		c.SecretKey,
	)
	essClient.EnableAsync(5, 10)
	if err != nil {
		return nil, err
	}

	return &Client{
		ess: essClient,
	}, nil
}

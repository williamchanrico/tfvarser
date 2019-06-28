package ecs

import (
	ecssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// Client is the ess client
type Client struct {
	*ecssdk.Client
}

// New returns a new ess client
func New(c *ecssdk.Client) *Client {
	return &Client{c}
}

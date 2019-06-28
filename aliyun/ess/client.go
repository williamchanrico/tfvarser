package ess

import (
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// Client is the ess client
type Client struct {
	*esssdk.Client
}

// New returns a new ess client
func New(c *esssdk.Client) *Client {
	return &Client{c}
}

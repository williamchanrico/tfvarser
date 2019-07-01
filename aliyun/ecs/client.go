package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// Client is the ess client
type Client struct {
	*ecs.Client
}

// New returns a new ess client
func New(c sdk.Client) *Client {
	return &Client{&ecs.Client{Client: c}}
}

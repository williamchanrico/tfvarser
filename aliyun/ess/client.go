package ess

import "github.com/aliyun/alibaba-cloud-sdk-go/sdk"
import "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"

// Client is the ess client
type Client struct {
	*ess.Client
}

// New returns a new ess client
func New(c sdk.Client) *Client {
	return &Client{&ess.Client{Client: c}}
}

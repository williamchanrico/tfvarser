package ess

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// LifecycleHook struct is mapped to lifecycle hook template
type LifecycleHook struct {
	LifecycleHookName   string
	LifecycleHookID     string
	LifecycleTransition string
	DefaultResult       string
	HeartbeatTimeout    int
}

// GetLifecycleHooks returns list of lifecyclehooks
func (c *Client) GetLifecycleHooks(scalingGroupID string) ([]LifecycleHook, error) {
	req := esssdk.CreateDescribeLifecycleHooksRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	lifecycleHooks := []LifecycleHook{}

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeLifecycleHooks(req)
		if err != nil {
			return nil, err
		}

		for _, lh := range resp.LifecycleHooks.LifecycleHook {
			lifecycleHook := LifecycleHook{}
			lifecycleHook.LifecycleHookName = lh.LifecycleHookName
			lifecycleHook.LifecycleHookID = lh.LifecycleHookId
			lifecycleHook.LifecycleTransition = lh.LifecycleTransition
			lifecycleHook.HeartbeatTimeout = lh.HeartbeatTimeout
			lifecycleHook.DefaultResult = lh.DefaultResult
			lifecycleHooks = append(lifecycleHooks, lifecycleHook)
		}
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.LifecycleHooks.LifecycleHook))
	}

	return lifecycleHooks, nil
}

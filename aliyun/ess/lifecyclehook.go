package ess

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// LifecycleHook struct is mapped to lifecycle hook template
type LifecycleHook struct {
	LifecycleHookName   string
	LifecycleHookID     string
	ScalingGroupName    string
	LifecycleTransition string
	DefaultResult       string
	HeartbeatTimeout    int
}

// GetLifecycleHooks returns list of lifecyclehooks
// scalingGroupName is only used to fill the struct, not for the request
func (c *Client) GetLifecycleHooks(scalingGroupID, scalingGroupName string) ([]LifecycleHook, error) {
	req := esssdk.CreateDescribeLifecycleHooksRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	lifecycleHooks := []LifecycleHook{}

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.ess.DescribeLifecycleHooks(req)
		if err != nil {
			return nil, err
		}

		for _, lh := range resp.LifecycleHooks.LifecycleHook {
			lifecycleHook := LifecycleHook{}
			lifecycleHook.LifecycleHookName = lh.LifecycleHookName
			lifecycleHook.LifecycleHookID = lh.LifecycleHookId
			lifecycleHook.ScalingGroupName = scalingGroupName // Needed for template
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

// DeleteLifecycleHook will delete the lifecycle hook
func (c *Client) DeleteLifecycleHook(lifecycleHookID string) error {
	req := esssdk.CreateDeleteLifecycleHookRequest()
	req.LifecycleHookId = lifecycleHookID

	_, err := c.ess.DeleteLifecycleHook(req)
	if err != nil {
		return err
	}

	return nil
}

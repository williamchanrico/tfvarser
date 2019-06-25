package ess

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// LifecycleHook struct is mapped to lifecycle hook template
type LifecycleHook struct {
	LifecycleHookName   string
	ScalingGroupName    string
	LifecycleTransition string
	DefaultResult       string
	HeartbeatTimeout    int
}

// LifecycleHookTmpl is the tfvars template
var LifecycleHookTmpl = `
terragrunt = {
  include {
    path = "${find_in_parent_folders()}"
  }
  terraform {
    source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-lifecycle-hook"
  }
}

# ESS scaling group
esssg_remote_state_bucket = "tkpd-tg-alicloud"
esssg_remote_state_key    = "{{ .ScalingGroupName }}/autoscale/ess-scaling-group/terraform.tfstate"

# MNS queue
mq_remote_state_bucket = "tkpd-tg-alicloud"
mq_remote_state_key    = "general/mns-queues/autoscaledown-event/terraform.tfstate"

# ESS lifecycle hook
esslh_name                 = "{{ .LifecycleHookName }}"
esslh_lifecycle_transition = "{{ .LifecycleTransition }}"
esslh_default_result       = "{{ .DefaultResult }}"
esslh_heartbeat_timeout    = {{ .HeartbeatTimeout }}
`

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

package ess

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// ScalingRule struct to map to the tfvars template
type ScalingRule struct {
	ScalingRuleName  string
	ScalingGroupID   string
	ScalingGroupName string
	AdjustmentType   string
	AdjustmentValue  int
	Cooldown         int
}

// ScalingRuleTmpl is the tfvars template
var ScalingRuleTmpl = `
terragrunt = {
  include {
    path = "${find_in_parent_folders()}"
  }
  terraform {
    source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-scaling-rule"
  }
}

# ESS scaling group
esssg_remote_state_bucket = "tkpd-tg-alicloud"
esssg_remote_state_key    = "{{ .ScalingGroupName }}/autoscale/ess-scaling-group/terraform.tfstate"

# ESS scaling rule
esssr_scaling_rule_name = "{{ .ScalingRuleName }}"
esssr_adjustment_type   = "{{ .AdjustmentType }}"
esssr_adjustment_value  = "{{ .AdjustmentValue }}"
esssr_cooldown          = {{ .Cooldown }}
`

// GetScalingRules returns list of scaling rule for the given scaling group
// scalingGroupName is only used to fill the struct, not for the request
func (c *Client) GetScalingRules(scalingGroupID, scalingGroupName string) ([]ScalingRule, error) {
	req := esssdk.CreateDescribeScalingRulesRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var scalingRules []ScalingRule

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.ess.DescribeScalingRules(req)
		if err != nil {
			return nil, err
		}

		for _, sr := range resp.ScalingRules.ScalingRule {
			scalingRule := ScalingRule{}
			scalingRule.ScalingRuleName = sr.ScalingRuleName
			scalingRule.ScalingGroupID = sr.ScalingGroupId
			scalingRule.ScalingGroupName = scalingGroupName // Needed for template
			scalingRule.AdjustmentType = sr.AdjustmentType
			scalingRule.AdjustmentValue = sr.AdjustmentValue
			scalingRule.Cooldown = sr.Cooldown

			scalingRules = append(scalingRules, scalingRule)
		}

		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.ScalingRules.ScalingRule))
	}

	return scalingRules, nil
}

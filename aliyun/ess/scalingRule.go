package ess

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// ScalingRule struct to map to the tfvars template
type ScalingRule struct {
	ScalingRuleName string
	ScalingRuleID   string
	ScalingGroupID  string
	AdjustmentType  string
	AdjustmentValue int
	Cooldown        int
}

// GetScalingRules returns list of scaling rule for the given scaling group
func (c *Client) GetScalingRules(scalingGroupID string) ([]ScalingRule, error) {
	req := esssdk.CreateDescribeScalingRulesRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var scalingRules []ScalingRule

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeScalingRules(req)
		if err != nil {
			return nil, err
		}

		for _, sr := range resp.ScalingRules.ScalingRule {
			scalingRule := ScalingRule{}
			scalingRule.ScalingRuleName = sr.ScalingRuleName
			scalingRule.ScalingRuleID = sr.ScalingRuleId
			scalingRule.ScalingGroupID = sr.ScalingGroupId
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

// GetScalingRuleByAri returns the scaling rule matched by it's ari
func (c *Client) GetScalingRuleByAri(ari string) (ScalingRule, error) {
	req := esssdk.CreateDescribeScalingRulesRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingRuleAri1 = ari

	resp, err := c.DescribeScalingRules(req)
	if err != nil {
		return ScalingRule{}, err
	}

	if len(resp.ScalingRules.ScalingRule) < 1 {
		return ScalingRule{}, errors.New("Scaling rule not found")
	}

	sr := resp.ScalingRules.ScalingRule[0]

	return ScalingRule{
		ScalingRuleName: sr.ScalingRuleName,
		ScalingRuleID:   sr.ScalingRuleId,
		ScalingGroupID:  sr.ScalingGroupId,
		AdjustmentType:  sr.AdjustmentType,
		AdjustmentValue: sr.AdjustmentValue,
		Cooldown:        sr.Cooldown,
	}, nil
}

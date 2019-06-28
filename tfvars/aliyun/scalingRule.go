package aliyun

import (
	"fmt"
	"io"
	"text/template"

	"github.com/williamchanrico/tfvarser/aliyun/ess"
)

const (
	// ScalingRuleKey is the reference key for scaling rule object
	ScalingRuleKey = "ess-scaling-rule"
)

// ScalingRule generator struct
type ScalingRule struct {
	ess.ScalingRule
	ScalingGroup ess.ScalingGroup
	ServiceName  string
}

// NewScalingRule return a generator for the scaling rule
func NewScalingRule(sr ess.ScalingRule, sg ess.ScalingGroup, serviceName string) *ScalingRule {
	return &ScalingRule{
		ScalingRule:  sr,
		ScalingGroup: sg,
		ServiceName:  serviceName,
	}
}

// Name returns the name of this tfvars generator
func (s *ScalingRule) Name() string {
	return fmt.Sprintf("%s-%s", s.Kind(), s.ScalingRuleName)
}

// Kind returns the key reference to this provider and object
func (s *ScalingRule) Kind() string {
	return fmt.Sprintf("%s-%s", Provider, ScalingRuleKey)
}

// Execute a scaling rule raw string
func (s *ScalingRule) Execute(w io.Writer, tmpl *template.Template) error {
	if err := tmpl.Execute(w, s); err != nil {
		return err
	}

	return nil
}

// Template returns the template
func (s *ScalingRule) Template() string {
	tmpl := `terragrunt = {
  include {
    path = "${find_in_parent_folders()}"
  }
  terraform {
    source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-scaling-rule"
  }
}

# ESS scaling group (ID: {{ .ScalingGroup.ScalingGroupID }})
esssg_remote_state_bucket = "tkpd-tg-alicloud"
esssg_remote_state_key    = "{{ .ServiceName }}/autoscale/ess-scaling-group/terraform.tfstate"

# ESS scaling rule
esssr_scaling_rule_name = "{{ .ScalingRuleName }}"
esssr_adjustment_type   = "{{ .AdjustmentType }}"
esssr_adjustment_value  = "{{ .AdjustmentValue }}"
esssr_cooldown          = {{ .Cooldown }}

# Import command
# terragrunt import alicloud_ess_scaling_rule.esssr {{ .ScalingRuleID }}
`

	return tmpl
}

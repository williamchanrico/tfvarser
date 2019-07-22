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

	Extras map[string]interface{}
}

// NewScalingRule return a generator for the scaling rule
func NewScalingRule(sr ess.ScalingRule, sg ess.ScalingGroup, extras map[string]interface{}) *ScalingRule {
	return &ScalingRule{
		ScalingRule:  sr,
		ScalingGroup: sg,
		Extras:       extras,
	}
}

// Name returns the name of this tfvars generator
func (s *ScalingRule) Name() string {
	return fmt.Sprintf("%s-%s-%s", s.Provider(), s.Kind(), s.ScalingRuleName)
}

// Kind returns the key reference to this object
func (s *ScalingRule) Kind() string {
	return fmt.Sprintf("%s", ScalingRuleKey)
}

// Provider returns the key reference to this provider
func (s *ScalingRule) Provider() string {
	return fmt.Sprintf("%s", Provider)
}

// Execute a scaling rule raw string
func (s *ScalingRule) Execute(w io.Writer, tmpl *template.Template) error {
	return tmpl.Execute(w, s)
}

// Template returns the template
// func (s *ScalingRule) Template() string {
//     tmpl := `include {
//   path = "${find_in_parent_folders()}"
// }

// terraform {
//   source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-scaling-rule"
// }

// inputs = {
//   # ESS Scaling Group
//   esssg_remote_state_bucket = "tkpd-tg-alicloud"
//   esssg_remote_state_key    = "{{ index .Extras "serviceName" }}/autoscale/ess-scaling-group/terraform.tfstate"

//   # ESS Scaling Rule
//   esssr_scaling_rule_name = "{{ trimPrefix .ScalingRuleName "tf-" }}"
//   esssr_adjustment_type   = "{{ .AdjustmentType }}"
//   esssr_adjustment_value  = "{{ .AdjustmentValue }}"
//   esssr_cooldown          = {{ .Cooldown }}
// }`

//     return tmpl
// }

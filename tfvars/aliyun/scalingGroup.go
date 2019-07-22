package aliyun

import (
	"fmt"
	"io"
	"text/template"

	"github.com/williamchanrico/tfvarser/aliyun/ess"
)

const (
	// ScalingGroupKey is the reference key for scaling group object
	ScalingGroupKey = "ess-scaling-group"
)

// ScalingGroup generator struct
type ScalingGroup struct {
	ess.ScalingGroup

	Extras map[string]interface{}
}

// NewScalingGroup return a generator for the scaling group
func NewScalingGroup(sg ess.ScalingGroup, extras map[string]interface{}) *ScalingGroup {
	return &ScalingGroup{
		ScalingGroup: sg,
		Extras:       extras,
	}
}

// Name returns the name of this tfvars generator
func (s *ScalingGroup) Name() string {
	return fmt.Sprintf("%s-%s-%s", s.Provider(), s.Kind(), s.ScalingGroupName)
}

// Kind returns the key reference to this object
func (s *ScalingGroup) Kind() string {
	return fmt.Sprintf("%s", ScalingGroupKey)
}

// Provider returns the key reference to this provider
func (s *ScalingGroup) Provider() string {
	return fmt.Sprintf("%s", Provider)
}

// Execute a scaling group raw string
func (s *ScalingGroup) Execute(w io.Writer, tmpl *template.Template) error {
	return tmpl.Execute(w, s)
}

// Template returns the template
// func (s *ScalingGroup) Template() string {
//     tmpl := `include {
//   path = "${find_in_parent_folders()}"
// }

// terraform {
//   source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-scaling-group"
// }

// inputs = {
//   # VPC VSwitch
//   vsw_remote_state_bucket = "tkpd-tg-alicloud-infra"
//   vsw_remote_state_keys   = [
//      "vswitches/app/terraform.tfstate",
//      "vswitches/app-2/terraform.tfstate"
//   ]

//   # ESS Scaling Group
//   esssg_name = "{{ trimPrefix .ScalingGroupName "tf-" }}"

//   esssg_min_size = {{ .MinSize }}
//   esssg_max_size = {{ .MaxSize }}

//   esssg_removal_policies = [
//   {{ range $index, $element := .RemovalPolicies }}{{- if $index }},
//   {{- end }}{{ if not $index }} {{ end }} "{{ $element -}}"{{ end }}
//   ]

//   esssg_multi_az_policy  = "{{ .MultiAZPolicy }}"
//   {{- if .LoadBalancerIDs }}
//   esssg_loadbalancer_ids = [
//   {{ range $index, $element := .LoadBalancerIDs }}{{- if $index }},
//   {{- end }}{{ if not $index }} {{ end }} "{{ $element -}}"{{ end }}
//   ]{{ end }}
// }`

//     return tmpl
// }

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
	svc ess.ScalingGroup
}

// NewScalingGroup return a generator for the scaling group
func NewScalingGroup(sg ess.ScalingGroup) *ScalingGroup {
	return &ScalingGroup{
		svc: sg,
	}
}

// Name returns the name of this tfvars generator
func (s *ScalingGroup) Name() string {
	return fmt.Sprintf("%s-%s", s.Kind(), s.svc.ScalingGroupName)
}

// Kind returns the key reference to this provider and object
func (s *ScalingGroup) Kind() string {
	return fmt.Sprintf("%s-%s", Provider, ScalingGroupKey)
}

// Execute a scaling group raw string
func (s *ScalingGroup) Execute(w io.Writer, tmpl *template.Template) error {
	if err := tmpl.Execute(w, s.svc); err != nil {
		return err
	}

	return nil
}

// Template returns the template
func (s *ScalingGroup) Template() string {
	tmpl := `terragrunt = {
  include {
    path = "${find_in_parent_folders()}"
  }
  terraform {
    source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-scaling-group"
  }
}

# Name of the scaling group (ID: {{ .ScalingGroupID }})
esssg_name = "{{ .ScalingGroupName }}"

# Minimum and maximum number of VMs in the scaling group
esssg_min_size = {{ .MinSize }}
esssg_max_size = {{ .MaxSize }}

# When downscaling, this specifies the order of VMs selected for removal
esssg_removal_policies = [
{{ range $index, $element := .RemovalPolicies }}{{- if $index }},
{{- end }}{{ if not $index }} {{ end }} "{{ $element -}}"{{ end }}
]

# VSwitches that will be used for created VMs, selection algorithm is based on esssg_multi_az_policy
esssg_vsw_ids          = [
{{ range $index, $element := .VSwitchIDs }}{{- if $index }},
{{- end }}{{ if not $index }} {{ end }} "{{ $element -}}"{{ end }}
]

# The order of VSwitches selected when creating new VMs
esssg_multi_az_policy  = "{{ .MultiAZPolicy }}"
`

	return tmpl
}

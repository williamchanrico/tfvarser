package aliyun

import (
	"fmt"
	"io"
	"text/template"

	"github.com/williamchanrico/tfvarser/aliyun/ess"
)

const (
	// ScalingConfigurationKey is the reference key for scaling configuration object
	ScalingConfigurationKey = "ess-scaling-configuration"
)

// ScalingConfiguration generator struct
type ScalingConfiguration struct {
	svc ess.ScalingConfiguration
}

// NewScalingConfiguration return a generator for the scaling configuration
func NewScalingConfiguration(sc ess.ScalingConfiguration) *ScalingConfiguration {
	return &ScalingConfiguration{
		svc: sc,
	}
}

// Name returns the name of this tfvars generator
func (s *ScalingConfiguration) Name() string {
	return fmt.Sprintf("%s-%s", s.Kind(), s.svc.ScalingConfigurationName)
}

// Kind returns the key reference to this provider and object
func (s *ScalingConfiguration) Kind() string {
	return fmt.Sprintf("%s-%s", Provider, ScalingConfigurationKey)
}

// Execute a scaling configuration raw string
func (s *ScalingConfiguration) Execute(w io.Writer, tmpl *template.Template) error {
	if err := tmpl.Execute(w, s.svc); err != nil {
		return err
	}

	return nil
}

// Template returns the template
func (s *ScalingConfiguration) Template() string {
	tmpl := `terragrunt = {
  include {
    path = "${find_in_parent_folders()}"
  }
  terraform {
    source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-scaling-configuration"
  }
}

# ESS scaling group
esssg_remote_state_bucket = "tkpd-tg-alicloud"
esssg_remote_state_key    = "{{ .ScalingGroupName }}/autoscale/ess-scaling-group/terraform.tfstate"

# Security group
sg_remote_state_bucket = "tkpd-tg-alicloud-infra"
sg_remote_state_key    = "security-groups/intranet/security-group/terraform.tfstate"

# ECS Images
images_name_regex = "{{ .ImageName }}"

# ESS scaling configuration
esssc_scaling_configuration_name = "{{ .ScalingGroupName }}"
esssc_instance_name              = "{{ .ScalingGroupName }}"
esssc_instance_types             = [
{{ range $index, $element := .InstanceTypes }}{{- if $index }},
{{- end }}{{- if not $index }} {{ end }} "{{ $element -}}"{{ end }}
]
esssc_enable                     = {{ .Enable }}
esssc_active                     = {{ .Active }}
esssc_key_name                   = "{{ .KeyPairName }}"
esssc_role_name                  = "{{ .RAMRoleName }}"

esssc_user_data = <<EOF
{{ .UserData }}
EOF

esssc_tags_tribe     = ""
esssc_tags_team      = ""
esssc_tags_hostgroup = "{{ index .Tags "hostgroup" }}"
esssc_tags_type      = ""

{{ if .Tags.consul_tags }}
esssc_optional_tags = {
  "consul_tags" = "{{ index .Tags "consul_tags" }}"
}
{{ end }}
`

	return tmpl
}

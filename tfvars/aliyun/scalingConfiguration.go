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
	ess.ScalingConfiguration
	ScalingGroup ess.ScalingGroup

	Extras map[string]interface{}
}

// NewScalingConfiguration return a generator for the scaling configuration
func NewScalingConfiguration(sc ess.ScalingConfiguration, sg ess.ScalingGroup, extras map[string]interface{}) *ScalingConfiguration {
	return &ScalingConfiguration{
		ScalingConfiguration: sc,
		ScalingGroup:         sg,
		Extras:               extras,
	}
}

// Name returns the name of this tfvars generator
func (s *ScalingConfiguration) Name() string {
	return fmt.Sprintf("%s-%s", s.Kind(), s.ScalingConfigurationName)
}

// Kind returns the key reference to this provider and object
func (s *ScalingConfiguration) Kind() string {
	return fmt.Sprintf("%s-%s", Provider, ScalingConfigurationKey)
}

// Execute a scaling configuration raw string
func (s *ScalingConfiguration) Execute(w io.Writer, tmpl *template.Template) error {
	return tmpl.Execute(w, s)
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

# ESS Scaling group (ID: {{ .ScalingGroup.ScalingGroupID }})
esssg_remote_state_bucket = "tkpd-tg-alicloud"
esssg_remote_state_key    = "{{ index .Extras "serviceName" }}/autoscale/ess-scaling-group/terraform.tfstate"

# ECS Security group
sg_remote_state_bucket = "tkpd-tg-alicloud-infra"
sg_remote_state_key    = "security-groups/intranet/security-group/terraform.tfstate"

# ECS Images
images_name_regex = "^{{ index .Extras "imageName" }}$"

# ESS Scaling configuration
esssc_scaling_configuration_name = "{{ trimPrefix .ScalingConfigurationName "tf-" }}"
esssc_instance_name              = "{{ index .Extras "serviceName" }}"
esssc_instance_types             = [
{{ range $index, $element := .InstanceTypes }}{{- if $index }},
{{- end }}{{- if not $index }} {{ end }} "{{ $element -}}"{{ end }}
]
esssc_enable                     = {{ .Enable }}
esssc_active                     = {{ .Active }}
esssc_key_name                   = "{{ .KeyPairName }}"
esssc_role_name                  = "{{ .RAMRoleName }}"

esssc_user_data = <<EOF
{{ .UserData -}}

EOF

esssc_tags_tribe     = "{{ if .Tags.tribe }}{{ index .Tags "tribe" }}{{ end }}"
esssc_tags_team      = "{{ if .Tags.team }}{{ index .Tags "team" }}{{ end }}"
esssc_tags_hostgroup = "{{ if .Tags.hostgroup }}{{ index .Tags "hostgroup" }}{{ end }}"
esssc_tags_type      = "{{ if .Tags.type }}{{ index .Tags "type" }}{{ end }}"
{{ if .Tags.consul_tags }}
esssc_optional_tags = {
  "consul_tags" = "{{ index .Tags "consul_tags" }}"
}
{{ end }}
# Import command
# terragrunt import alicloud_ess_scaling_configuration.esssc {{ .ScalingConfigurationID }}
`

	return tmpl
}

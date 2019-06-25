package ess

import (
	"encoding/base64"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// ScalingConfiguration struct to map to the tfvars template
type ScalingConfiguration struct {
	ScalingConfigurationName string
	ScalingGroupID           string
	ScalingGroupName         string
	ImageID                  string
	ImageName                string
	InstanceName             string
	InstanceType             string
	InstanceTypes            []string
	Enable                   bool
	Active                   bool
	KeyPairName              string
	RAMRoleName              string
	UserData                 string
	Tags                     map[string]string
}

// ScalingConfigurationTmpl is the tfvars template
var ScalingConfigurationTmpl = `
terragrunt = {
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

# ESS scaling configuration
esssc_scaling_configuration_name = "{{ .ScalingConfigurationName }}"
esssc_image_id                   = "{{ .ImageID }}"
esssc_image_name                 = "{{ .ImageName }}"
esssc_instance_name              = "{{ .InstanceName }}"
esssc_instance_type              = "{{ .InstanceType }}"
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

// GetScalingConfigurations returns list of scaling rule for the given scaling group
// scalingGroupName is only used to fill the struct, not for the request
func (c *Client) GetScalingConfigurations(scalingGroupID, scalingGroupName string) ([]ScalingConfiguration, error) {
	req := esssdk.CreateDescribeScalingConfigurationsRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var scalingConfigurations []ScalingConfiguration

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.ess.DescribeScalingConfigurations(req)
		if err != nil {
			return nil, err
		}

		for _, sc := range resp.ScalingConfigurations.ScalingConfiguration {
			scalingConfiguration := ScalingConfiguration{}
			scalingConfiguration.ScalingConfigurationName = sc.ScalingConfigurationName
			scalingConfiguration.ScalingGroupID = sc.ScalingGroupId
			scalingConfiguration.ScalingGroupName = scalingGroupName // Needed for template
			scalingConfiguration.ImageID = sc.ImageId
			scalingConfiguration.ImageName = sc.ImageName
			scalingConfiguration.InstanceName = sc.InstanceName
			scalingConfiguration.InstanceType = sc.InstanceType
			scalingConfiguration.InstanceTypes = sc.InstanceTypes.InstanceType
			scalingConfiguration.Enable = false
			scalingConfiguration.Active = false
			if sc.LifecycleState == "Active" {
				scalingConfiguration.Active = true
				scalingConfiguration.Enable = true
			}
			scalingConfiguration.KeyPairName = sc.KeyPairName
			scalingConfiguration.RAMRoleName = sc.RamRoleName

			var ok bool
			scalingConfiguration.UserData, ok = base64Decode(sc.UserData)
			if !ok {
				scalingConfiguration.UserData = ""
			}

			scalingConfiguration.Tags = make(map[string]string)
			for _, tag := range sc.Tags.Tag {
				scalingConfiguration.Tags[strings.ToLower(tag.Key)] = tag.Value
			}

			scalingConfigurations = append(scalingConfigurations, scalingConfiguration)
		}

		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.ScalingConfigurations.ScalingConfiguration))
	}

	return scalingConfigurations, nil
}

func base64Decode(str string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", false
	}

	return string(data), true
}

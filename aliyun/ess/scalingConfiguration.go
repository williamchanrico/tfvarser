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
	ScalingConfigurationID   string
	ScalingGroupID           string
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

// GetScalingConfigurations returns list of scaling rule for the given scaling group
func (c *Client) GetScalingConfigurations(scalingGroupID string) ([]ScalingConfiguration, error) {
	req := esssdk.CreateDescribeScalingConfigurationsRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var scalingConfigurations []ScalingConfiguration

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeScalingConfigurations(req)
		if err != nil {
			return nil, err
		}

		for _, sc := range resp.ScalingConfigurations.ScalingConfiguration {
			scalingConfiguration := ScalingConfiguration{}
			scalingConfiguration.ScalingConfigurationName = sc.ScalingConfigurationName
			scalingConfiguration.ScalingConfigurationID = sc.ScalingConfigurationId
			scalingConfiguration.ScalingGroupID = sc.ScalingGroupId
			scalingConfiguration.ImageID = sc.ImageId
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

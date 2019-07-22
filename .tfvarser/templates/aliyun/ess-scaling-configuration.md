# ess-scaling-configuration

## Available Fields

```
.ScalingConfiguration
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

.ScalingGroup
	< see ess-scaling-group.md >

.Extras
	map["serviceName"] # e.g. "{{ index .Extras "serviceName" }}"
	map["imageName"]
```

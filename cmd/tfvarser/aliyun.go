package tfvarser

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/williamchanrico/tfvarser/aliyun"
	"github.com/williamchanrico/tfvarser/tfvars"
	tfvarsaliyun "github.com/williamchanrico/tfvarser/tfvars/aliyun"
)

func aliyunProvider(appFlags *Flags, cfg Config) (int, error) {
	switch appFlags.ProviderObj {
	case "ess":
		return aliyunAutoscaleObjects(appFlags, cfg)

	default:
		return 1, errors.New("Object is not supported")
	}
}

// aliyunAutoscaleObjects generates autoscale related objects
// generated structure:
// ├── testapp
// │   ├── autoscale
// │   │   ├── ess-alarms
// │   │   │   ├── go-testapp-downscale
// │   │   │   │   └── terraform.tfvars
// │   │   │   └── go-testapp-upscale
// │   │   │       └── terraform.tfvars
// │   │   ├── ess-lifecycle-hooks
// │   │   │   ├── autoscaledown-event-mns-queue
// │   │   │   │   └── terraform.tfvars
// │   │   │   └── autoscaleup-event-mns-queue
// │   │   │       └── terraform.tfvars
// │   │   ├── ess-scaling-configurations
// │   │   │   ├── go-testapp-1c-1gb
// │   │   │   │   └── terraform.tfvars
// │   │   │   └── go-testapp-1c-500mb
// │   │   │       └── terraform.tfvars
// │   │   ├── ess-scaling-group
// │   │   │   └── terraform.tfvars
// │   │   └── ess-scaling-rules
// │   │       ├── auto-downscale
// │   │       │   └── terraform.tfvars
// │   │       └── auto-upscale
// │   │           └── terraform.tfvars
//
func aliyunAutoscaleObjects(appFlags *Flags, cfg Config) (int, error) {
	aliClient, err := aliyun.New(aliyun.Config{
		AccessKey: cfg.AlicloudAccessKey,
		SecretKey: cfg.AlicloudSecretKey,
		RegionID:  cfg.AlicloudRegionID,
	})
	if err != nil {
		return 1, err
	}

	scalingGroups, err := aliClient.ESS.GetScalingGroups(context.Background())
	if err != nil {
		return 1, err
	}

	limitNames := strings.Split(appFlags.LimitNames, ",")
	limitIDs := strings.Split(appFlags.LimitIDs, ",")

	for _, sg := range scalingGroups {
		if !(contains(limitNames, sg.ScalingGroupName) || contains(limitIDs, sg.ScalingGroupID)) {
			continue
		}
		fmt.Printf("Generating tfvars for scaling group: %v\n", sg.ScalingGroupName)

		serviceName := parseServiceNameFromScalingGroup(sg.ScalingGroupName)
		serviceDir := path.Join(".", serviceName, "autoscale")

		// Scaling Group
		scalingGroupDir := path.Join(serviceDir, "ess-scaling-group")
		sgGenerator := tfvars.New(tfvarsaliyun.NewScalingGroup(sg))
		sgGenerator.Generate(scalingGroupDir, "terraform.tfvars")

		// Scaling Rule
		scalingRules, err := aliClient.ESS.GetScalingRules(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, err
		}
		scalingRuleParentDir := path.Join(serviceDir, "ess-scaling-rules")
		for _, sr := range scalingRules {
			scalingRuleDir := path.Join(scalingRuleParentDir, sr.ScalingRuleName)
			srGenerator := tfvars.New(tfvarsaliyun.NewScalingRule(sr))
			srGenerator.Generate(scalingRuleDir, "terraform.tfvars")
		}

		// Alarm or Event-trigger task
		alarms, err := aliClient.ESS.GetAlarms(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}
		alarmParentDir := path.Join(serviceDir, "ess-alarms")
		for _, al := range alarms {
			alarmDir := path.Join(alarmParentDir, al.AlarmName)
			alGenerator := tfvars.New(tfvarsaliyun.NewAlarm(al))
			alGenerator.Generate(alarmDir, "terraform.tfvars")
		}

		// // Lifecycle Hook
		lifecycleHooks, err := aliClient.ESS.GetLifecycleHooks(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}
		lifecycleHookParentDir := path.Join(serviceDir, "ess-lifecycle-hooks")
		for _, lh := range lifecycleHooks {
			lifecycleHookDir := path.Join(lifecycleHookParentDir, lh.LifecycleHookName)
			lhGenerator := tfvars.New(tfvarsaliyun.NewLifecycleHook(lh))
			lhGenerator.Generate(lifecycleHookDir, "terraform.tfvars")
		}

		// Scaling Configuration
		scalingConfigurations, err := aliClient.ESS.GetScalingConfigurations(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}
		scalingConfigurationParentDir := path.Join(serviceDir, "ess-scaling-configurations")
		for _, sc := range scalingConfigurations {
			// Template needs ImageName from ECS API
			// if imageName is empty, it's fine for now
			imageName, err := aliClient.ECS.GetImageNameByID(sc.ImageID)
			if err == nil {
				sc.ImageName = imageName
			}

			scalingConfigurationDir := path.Join(scalingConfigurationParentDir, sc.ScalingGroupName)
			scGenerator := tfvars.New(tfvarsaliyun.NewScalingConfiguration(sc))
			scGenerator.Generate(scalingConfigurationDir, "terraform.tfvars")
		}
	}

	return 0, nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// parseServiceNameFromScalingGroup to get service name without any scaling group specific tags
// this is purely optional out of current setup
func parseServiceNameFromScalingGroup(scalingGroupName string) string {
	ret := strings.TrimPrefix(scalingGroupName, "tf-go-")
	ret = strings.TrimPrefix(ret, "go-")
	ret = strings.TrimPrefix(ret, "node-")
	return ret
}

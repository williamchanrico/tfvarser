package tfvarser

import (
	"context"
	"fmt"
	"path"
	"strings"
	"text/template"
	"time"

	col "github.com/logrusorgru/aurora"

	"github.com/williamchanrico/tfvarser/aliyun"
	"github.com/williamchanrico/tfvarser/aliyun/ess"
	"github.com/williamchanrico/tfvarser/log"
	"github.com/williamchanrico/tfvarser/tfvars"
	tfvarsaliyun "github.com/williamchanrico/tfvarser/tfvars/aliyun"
)

func aliyunProvider(appFlags *Flags, cfg Config) (int, error) {
	switch appFlags.ProviderObj {
	case "ess":
		return aliyunAutoscaleObjects(appFlags, cfg)

	default:
		return 1, ErrObjNotSupported
	}
}

// aliyunAutoscaleObjects generates autoscale related objects
func aliyunAutoscaleObjects(appFlags *Flags, cfg Config) (int, error) {
	aliClient, err := aliyun.New(aliyun.Config{
		AccessKey: cfg.AlicloudAccessKey,
		SecretKey: cfg.AlicloudSecretKey,
		RegionID:  cfg.AlicloudRegionID,
	})
	if err != nil {
		return 1, err
	}

	funcMap := template.FuncMap{
		"trimPrefix": trimPrefix,
	}

	fmt.Printf("Querying %v from cloud provider\n", col.Cyan("Scaling Group(s)"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scalingGroups, err := aliClient.ESS.GetScalingGroupsWithAsync(ctx)
	if err != nil {
		return 1, err
	}
	log.Debugf("Retrieved %v Scaling Group(s)", len(scalingGroups))

	// We want to act on limit flags, so start processing them here
	// limitNames: will filter on ScalingGroupName
	// limitIDs: will filter on ScalingGroupID
	limitNames := strings.Split(appFlags.LimitNames, ",")
	if len(limitNames) > 0 {
		log.Debugf("Limiting search to Scaling Group with names: %v", limitNames)
	}
	limitIDs := strings.Split(appFlags.LimitIDs, ",")
	if len(limitIDs) > 0 {
		log.Debugf("Limiting search to Scaling Group with IDs: %v", limitIDs)
	}

	for _, sg := range scalingGroups {
		// Only process the ones in limit variables (if either of the limit flags were specified)
		if !(contains(limitNames, sg.ScalingGroupName) || contains(limitIDs, sg.ScalingGroupID)) {
			continue
		}
		fmt.Printf("\nGenerating tfvars for scaling group: %v\n", col.Green(sg.ScalingGroupName))

		// We want to separate every scaling group by service name
		// we will inject this service name to generators that need this service name
		extras := make(map[string]interface{})
		extras["serviceName"] = parseServiceNameFromScalingGroup(sg.ScalingGroupName)
		serviceDir := path.Join(".", extras["serviceName"].(string), "autoscale")

		// Scaling Group
		scalingGroupDir := path.Join(serviceDir, "ess-scaling-group")
		sgGenerator := tfvars.New(tfvarsaliyun.NewScalingGroup(sg, extras), funcMap)
		log.Debugf("Generating %v", path.Join(scalingGroupDir, "terraform.tfvars"))
		err = sgGenerator.Generate(scalingGroupDir, "terraform.tfvars")
		if err != nil {
			log.Errorf("Error generating %v: %v\n", path.Join(scalingGroupDir, "terraform.tfvars"), err.Error())
		}

		// Scaling Rule
		log.Debugf("Getting Scaling Rule(s) in %v", sg.ScalingGroupName)
		scalingRules, err := aliClient.ESS.GetScalingRules(sg.ScalingGroupID)
		if err != nil {
			return 1, err
		}
		scalingRuleParentDir := path.Join(serviceDir, "ess-scaling-rules")
		for _, sr := range scalingRules {
			// Replace scaling rule name with auto-{upscale/downscale} instead
			// when matched a criteria
			if strings.Contains(sr.ScalingRuleName, "downscale") {
				sr.ScalingRuleName = "auto-downscale"
			} else if strings.Contains(sr.ScalingRuleName, "upscale") {
				sr.ScalingRuleName = "auto-upscale"
			}

			scalingRuleDir := path.Join(scalingRuleParentDir, strings.TrimPrefix(sr.ScalingRuleName, "tf-"))
			srGenerator := tfvars.New(tfvarsaliyun.NewScalingRule(sr, sg, extras), funcMap)
			log.Debugf("Generating %v", path.Join(scalingRuleDir, "terraform.tfvars"))
			err = srGenerator.Generate(scalingRuleDir, "terraform.tfvars")
			if err != nil {
				log.Errorf("Error generating %v: %v\n", path.Join(scalingRuleDir, "terraform.tfvars"), err.Error())
			}
		}

		// Alarm (Event-trigger task)
		log.Debugf("Getting Alarms related to %v", sg.ScalingGroupName)
		alarms, err := aliClient.ESS.GetAlarms(sg.ScalingGroupID)
		if err != nil {
			return 1, nil
		}
		alarmParentDir := path.Join(serviceDir, "ess-alarms")
		for _, al := range alarms {
			// An alarm may have scaling rules, and we want them if they exist
			sr := ess.ScalingRule{}
			if len(al.AlarmActions) > 0 {
				// but currently we only care about 1 scaling rule per alarm
				sr, _ = aliClient.ESS.GetScalingRuleByAri(al.AlarmActions[0])

				// Replace scaling rule name with auto-{upscale/downscale} when matched a criteria
				if strings.Contains(sr.ScalingRuleName, "-upscale") {
					sr.ScalingRuleName = "auto-upscale"
				} else if strings.Contains(sr.ScalingRuleName, "-downscale") {
					sr.ScalingRuleName = "auto-downscale"
				}
			}

			alarmDir := path.Join(alarmParentDir, strings.TrimPrefix(al.AlarmName, "tf-"))
			alGenerator := tfvars.New(tfvarsaliyun.NewAlarm(al, sg, sr, extras), funcMap)
			log.Debugf("Generating %v", path.Join(alarmDir, "terraform.tfvars"))
			err = alGenerator.Generate(alarmDir, "terraform.tfvars")
			if err != nil {
				log.Errorf("Error generating %v: %v\n", path.Join(alarmDir, "terraform.tfvars"), err.Error())
			}
		}

		// Lifecycle Hook
		log.Debugf("Getting LifecycleHook(s) in %v", sg.ScalingGroupName)
		lifecycleHooks, err := aliClient.ESS.GetLifecycleHooks(sg.ScalingGroupID)
		if err != nil {
			return 1, nil
		}
		lifecycleHookParentDir := path.Join(serviceDir, "ess-lifecycle-hooks")
		for _, lh := range lifecycleHooks {
			// Use predefined name for LH, autoscale{up/down}-event-mns-queue
			if lh.LifecycleTransition == "SCALE_IN" {
				lh.LifecycleHookName = "autoscaledown-event-mns-queue"
			} else if lh.LifecycleTransition == "SCALE_OUT" {
				lh.LifecycleHookName = "autoscaleup-event-mns-queue"
			}

			lifecycleHookDir := path.Join(lifecycleHookParentDir, lh.LifecycleHookName)
			lhGenerator := tfvars.New(tfvarsaliyun.NewLifecycleHook(lh, sg, extras), funcMap)
			log.Debugf("Generating %v", path.Join(lifecycleHookDir, "terraform.tfvars"))
			err = lhGenerator.Generate(lifecycleHookDir, "terraform.tfvars")
			if err != nil {
				log.Errorf("Error generating %v: %v\n", path.Join(lifecycleHookDir, "terraform.tfvars"), err.Error())
			}
		}

		// Scaling Configuration
		log.Debugf("Getting Scaling Configuration(s) in %v", sg.ScalingGroupName)
		scalingConfigurations, err := aliClient.ESS.GetScalingConfigurations(sg.ScalingGroupID)
		if err != nil {
			return 1, nil
		}
		scalingConfigurationParentDir := path.Join(serviceDir, "ess-scaling-configurations")
		for _, sc := range scalingConfigurations {
			// Template needs ImageName from ECS API, will inject this to the generator
			extras["imageName"], err = aliClient.ECS.GetImageNameByID(sc.ImageID)
			if err != nil {
				extras["imageName"] = "IMAGE_NOT_FOUND_REPLACE_ME"
			}

			scalingConfigurationDir := path.Join(scalingConfigurationParentDir, strings.TrimPrefix(sc.ScalingConfigurationName, "tf-"))
			scGenerator := tfvars.New(tfvarsaliyun.NewScalingConfiguration(sc, sg, extras), funcMap)
			log.Debugf("Generating %v", path.Join(scalingConfigurationDir, "terraform.tfvars"))
			err = scGenerator.Generate(scalingConfigurationDir, "terraform.tfvars")
			if err != nil {
				log.Errorf("Error generating %v: %v\n", path.Join(scalingConfigurationDir, "terraform.tfvars"), err.Error())
			}
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

func trimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// parseServiceNameFromScalingGroup to get service name without any scaling group specific tags
// this is purely optional out of current setup
func parseServiceNameFromScalingGroup(scalingGroupName string) string {
	ret := strings.TrimPrefix(scalingGroupName, "tf-go-")
	ret = strings.TrimPrefix(ret, "go-")
	ret = strings.TrimPrefix(ret, "node-")
	return ret
}

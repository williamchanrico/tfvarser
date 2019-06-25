package tfvarser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/williamchanrico/tfvarser/aliyun/ess"
)

func aliyunProvider(appFlags *Flags, cfg Config) (int, error) {
	switch appFlags.ProviderObj {
	case "ess":
		return aliyunESSObj(appFlags, cfg)

	default:
		return 1, errors.New("Object is not supported")
	}
}

func aliyunESSObj(appFlags *Flags, cfg Config) (int, error) {
	essclient, err := ess.New(&ess.Config{
		AccessKey: cfg.AlicloudAccessKey,
		SecretKey: cfg.AlicloudSecretKey,
		RegionID:  cfg.AlicloudRegionID,
	})
	if err != nil {
		return 1, err
	}

	scalingGroups, err := essclient.GetScalingGroups(context.Background())
	if err != nil {
		return 1, err
	}

	tScalingGroup := template.Must(template.New("scalingGroup").Parse(ess.ScalingGroupTmpl))
	tScalingRule := template.Must(template.New("scalingRule").Parse(ess.ScalingRuleTmpl))
	tAlarm := template.Must(template.New("alarm").Parse(ess.AlarmTmpl))
	tLifecycleHook := template.Must(template.New("lifecycleHook").Parse(ess.LifecycleHookTmpl))
	tScalingConfiguration := template.Must(template.New("scalingConfiguration").Parse(ess.ScalingConfigurationTmpl))

	for _, sg := range scalingGroups {
		fmt.Printf("Generating scaling group: %v\n", sg.ScalingGroupName)

		serviceDir := path.Join(".", sg.ScalingGroupName, "autoscale")
		makeDirIfNotExists(serviceDir)

		scalingGroupDir := path.Join(serviceDir, "ess-scaling-group")
		makeDirIfNotExists(scalingGroupDir)
		f, err := os.Create(path.Join(scalingGroupDir, "terraform.tfvars"))
		if err != nil {
			fmt.Printf("Error %v: %v\n", sg.ScalingGroupName, err)
			continue
		}
		tScalingGroup.Execute(f, sg)
		f.Close()

		scalingRules, err := essclient.GetScalingRules(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, err
		}

		for _, sr := range scalingRules {
			scalingRuleParentDir := path.Join(serviceDir, "ess-scaling-rules")
			makeDirIfNotExists(scalingRuleParentDir)

			scalingRuleDir := path.Join(scalingRuleParentDir, sr.ScalingRuleName)
			makeDirIfNotExists(scalingRuleDir)
			f, err := os.Create(path.Join(scalingRuleDir, "terraform.tfvars"))
			if err != nil {
				fmt.Printf("Error %v: %v\n", sg.ScalingGroupName, err)
				continue
			}
			tScalingRule.Execute(f, sr)
			f.Close()
		}

		alarms, err := essclient.GetAlarms(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}

		alarmParentDir := path.Join(serviceDir, "ess-alarms")
		makeDirIfNotExists(alarmParentDir)
		for _, al := range alarms {
			alarmDir := path.Join(alarmParentDir, al.AlarmName)
			makeDirIfNotExists(alarmDir)
			f, err := os.Create(path.Join(alarmDir, "terraform.tfvars"))
			if err != nil {
				fmt.Printf("Error %v: %v\n", sg.ScalingGroupName, err)
				continue
			}

			tAlarm.Execute(f, al)
			f.Close()
		}

		lifecycleHooks, err := essclient.GetLifecycleHooks(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}

		lifecycleHookParentDir := path.Join(serviceDir, "ess-lifecycle-hooks")
		makeDirIfNotExists(lifecycleHookParentDir)
		for _, lh := range lifecycleHooks {
			lifecycleHookDir := path.Join(lifecycleHookParentDir, lh.LifecycleHookName)
			makeDirIfNotExists(lifecycleHookDir)
			f, err := os.Create(path.Join(lifecycleHookDir, "terraform.tfvars"))
			if err != nil {
				fmt.Printf("Error %v: %v\n", sg.ScalingGroupName, err)
				continue
			}

			tLifecycleHook.Execute(f, lh)
			f.Close()
		}

		scalingConfigurations, err := essclient.GetScalingConfigurations(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}

		scalingConfigurationParentDir := path.Join(serviceDir, "ess-scaling-configurations")
		makeDirIfNotExists(scalingConfigurationParentDir)
		for _, sc := range scalingConfigurations {
			scalingConfigurationDir := path.Join(scalingConfigurationParentDir, sc.ScalingConfigurationName)
			makeDirIfNotExists(scalingConfigurationDir)
			f, err := os.Create(path.Join(scalingConfigurationDir, "terraform.tfvars"))
			if err != nil {
				fmt.Printf("Error %v: %v\n", sg.ScalingGroupName, err)
				continue
			}

			tScalingConfiguration.Execute(f, sc)
			f.Close()
		}
		break
	}

	return 0, nil
}

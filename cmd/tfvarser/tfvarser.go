package tfvarser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/williamchanrico/tfvarser/aliyun/ess"
)

// Flags contains run-time flags
type Flags struct {
	Provider    string
	ProviderObj string
}

// Config contains tfvars config
type Config struct {
	AlicloudAccessKey string `envconfig:"ALICLOUD_ACCESS_KEY" required:"true"`
	AlicloudSecretKey string `envconfig:"ALICLOUD_SECRET_KEY" required:"true"`
	AlicloudRegionID  string `default:"ap-southeast-1"`
}

// Run is the entrypoint for main autoscaleapp process
func Run(appFlags *Flags, cfg Config) (int, error) {
	switch appFlags.Provider {
	case "ali":
		return aliyunProvider(appFlags, cfg)
	default:
		return 1, errors.New("Provider is not supported")
	}
}

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
		tScalingGroup.Execute(os.Stdout, sg)

		fmt.Printf("\n\n=============\n\n")
		fmt.Printf("\n\n=============\n\n")

		scalingRules, err := essclient.GetScalingRules(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, err
		}

		fmt.Println("Scaling Rule")
		for _, sr := range scalingRules {
			tScalingRule.Execute(os.Stdout, sr)
			fmt.Printf("\n\n=============\n\n")
		}
		alarms, err := essclient.GetAlarms(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}

		fmt.Println("Alarms")
		for _, al := range alarms {
			tAlarm.Execute(os.Stdout, al)
			fmt.Printf("\n\n=============\n\n")
		}

		lifecycleHooks, err := essclient.GetLifecycleHooks(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}
		fmt.Println("LifecycleHook")
		for _, lh := range lifecycleHooks {
			tLifecycleHook.Execute(os.Stdout, lh)
			fmt.Printf("\n\n=============\n\n")
		}

		scalingConfigurations, err := essclient.GetScalingConfigurations(sg.ScalingGroupID, sg.ScalingGroupName)
		if err != nil {
			return 1, nil
		}
		fmt.Println("Scaling Configuration")
		for _, sc := range scalingConfigurations {
			tScalingConfiguration.Execute(os.Stdout, sc)
			fmt.Printf("\n\n=============\n\n")
		}
		break
	}

	return 0, nil
}

package aliyun

import (
	"fmt"
	"io"
	"text/template"

	"github.com/williamchanrico/tfvarser/aliyun/ess"
)

const (
	// LifecycleHookKey is the reference key for lifecycle hook object
	LifecycleHookKey = "ess-lifecycle-hook"
)

// LifecycleHook generator struct
type LifecycleHook struct {
	ess.LifecycleHook
	ScalingGroup ess.ScalingGroup

	Extras map[string]interface{}
}

// NewLifecycleHook return a generator for the lifecycle hook
func NewLifecycleHook(lh ess.LifecycleHook, sg ess.ScalingGroup, extras map[string]interface{}) *LifecycleHook {
	return &LifecycleHook{
		LifecycleHook: lh,
		ScalingGroup:  sg,
		Extras:        extras,
	}
}

// Name returns the name of this tfvars generator
func (s *LifecycleHook) Name() string {
	return fmt.Sprintf("%s-%s", s.Kind(), s.LifecycleHookName)
}

// Kind returns the key reference to this provider and object
func (s *LifecycleHook) Kind() string {
	return fmt.Sprintf("%s-%s", Provider, LifecycleHookKey)
}

// Execute a lifecycle hook raw string
func (s *LifecycleHook) Execute(w io.Writer, tmpl *template.Template) error {
	return tmpl.Execute(w, s)
}

// Template returns the template
func (s *LifecycleHook) Template() string {
	tmpl := `include {
  path = "${find_in_parent_folders()}"
}

terraform {
  source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-lifecycle-hook"
}

inputs = {
  # ESS Scaling Group (ID: {{ .ScalingGroup.ScalingGroupID }})
  esssg_remote_state_bucket = "tkpd-tg-alicloud"
  esssg_remote_state_key    = "{{ index .Extras "serviceName" }}/autoscale/ess-scaling-group/terraform.tfstate"

  # MNS Queue
  mq_remote_state_bucket = "tkpd-tg-alicloud"
  mq_remote_state_key    = "general/mns-queues/{{ if eq .LifecycleTransition "SCALE_IN" }}autoscaledown-event{{ else }}autoscaleup-event{{ end }}/terraform.tfstate"

  # ESS Lifecycle Hook
  esslh_name                 = "{{ .LifecycleHookName }}"
  esslh_lifecycle_transition = "{{ .LifecycleTransition }}"
  esslh_default_result       = "{{ .DefaultResult }}"
  esslh_heartbeat_timeout    = {{ .HeartbeatTimeout }}
}

# Import command
# terragrunt import alicloud_ess_lifecycle_hook.esslh {{ .LifecycleHookID }}`

	return tmpl
}

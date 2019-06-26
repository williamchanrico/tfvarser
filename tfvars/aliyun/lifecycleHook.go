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
	svc ess.LifecycleHook
}

// NewLifecycleHook return a generator for the lifecycle hook
func NewLifecycleHook(lh ess.LifecycleHook) *LifecycleHook {
	return &LifecycleHook{
		svc: lh,
	}
}

// Name returns the name of this tfvars generator
func (s *LifecycleHook) Name() string {
	return fmt.Sprintf("%s-%s", s.Kind(), s.svc.LifecycleHookName)
}

// Kind returns the key reference to this provider and object
func (s *LifecycleHook) Kind() string {
	return fmt.Sprintf("%s-%s", Provider, LifecycleHookKey)
}

// Execute a lifecycle hook raw string
func (s *LifecycleHook) Execute(w io.Writer, tmpl *template.Template) error {
	if err := tmpl.Execute(w, s.svc); err != nil {
		return err
	}

	return nil
}

// Template returns the template
func (s *LifecycleHook) Template() string {
	tmpl := `terragrunt = {
  include {
    path = "${find_in_parent_folders()}"
  }
  terraform {
    source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-lifecycle-hook"
  }
}

# ESS scaling group
esssg_remote_state_bucket = "tkpd-tg-alicloud"
esssg_remote_state_key    = "{{ .ScalingGroupName }}/autoscale/ess-scaling-group/terraform.tfstate"

# MNS queue
mq_remote_state_bucket = "tkpd-tg-alicloud"
mq_remote_state_key    = "general/mns-queues/autoscaledown-event/terraform.tfstate"

# ESS lifecycle hook
esslh_name                 = "{{ if eq .LifecycleTransition "SCALE_IN" }}autoscaledown{{ else }}autoscaleup{{ end }}-event-mns-queue"
esslh_lifecycle_transition = "{{ .LifecycleTransition }}"
esslh_default_result       = "{{ .DefaultResult }}"
esslh_heartbeat_timeout    = {{ .HeartbeatTimeout }}
`

	return tmpl
}

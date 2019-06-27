package aliyun

import (
	"fmt"
	"io"
	"text/template"

	"github.com/williamchanrico/tfvarser/aliyun/ess"
)

const (
	// AlarmKey is the reference key for alarm object
	AlarmKey = "ess-alarm"
)

// Alarm generator struct
type Alarm struct {
	svc ess.Alarm
}

// NewAlarm return a generator for the alarm
func NewAlarm(al ess.Alarm) *Alarm {
	return &Alarm{
		svc: al,
	}
}

// Name returns the name of this tfvars generator
func (s *Alarm) Name() string {
	return fmt.Sprintf("%s-%s", s.Kind(), s.svc.AlarmName)
}

// Kind returns the key reference to this provider and object
func (s *Alarm) Kind() string {
	return fmt.Sprintf("%s-%s", Provider, AlarmKey)
}

// Execute a alarm raw string
func (s *Alarm) Execute(w io.Writer, tmpl *template.Template) error {
	if err := tmpl.Execute(w, s.svc); err != nil {
		return err
	}

	return nil
}

// Template returns the template
func (s *Alarm) Template() string {
	tmpl := `terragrunt = {
  include {
    path = "${find_in_parent_folders()}"
  }
  terraform {
    source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-alarm"
  }
}

# ESS scaling group (ID: {{ .ScalingGroupID }})
esssg_remote_state_bucket = "tkpd-tg-alicloud"
esssg_remote_state_key    = "{{ .ScalingGroupName }}/autoscale/ess-scaling-group/terraform.tfstate"

# ESS scaling rule
esssr_remote_state_bucket = "tkpd-tg-alicloud"
esssr_remote_state_key    = "{{ .ScalingGroupName }}/autoscale/ess-scaling-rules/{{ .ScalingRuleName }}/terraform.tfstate"

# ESS alarm
essa_name                = "{{ .AlarmName }}"
essa_metric_type         = "{{ .MetricType }}"
essa_metric_name         = "{{ .MetricName }}"
essa_period              = {{ .Period }}
essa_statistics          = "{{ .Statistics }}"
essa_comparison_operator = "{{ .ComparisonOperator }}"
essa_threshold           = {{ .Threshold }}
essa_evaluation_count    = {{ .EvaluationCount }}

# Import command
# terragrunt import alicloud_ess_alarm.essa {{ .AlarmID }}
`

	return tmpl
}

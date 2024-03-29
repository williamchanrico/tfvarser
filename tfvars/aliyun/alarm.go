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
	ess.Alarm
	ScalingRule  ess.ScalingRule
	ScalingGroup ess.ScalingGroup

	Extras map[string]interface{}
}

// NewAlarm return a generator for the alarm
func NewAlarm(al ess.Alarm, sg ess.ScalingGroup, sr ess.ScalingRule, extras map[string]interface{}) *Alarm {
	return &Alarm{
		Alarm:        al,
		ScalingRule:  sr,
		ScalingGroup: sg,
		Extras:       extras,
	}
}

// Name returns the name of this tfvars generator
func (s *Alarm) Name() string {
	return fmt.Sprintf("%s-%s-%s", s.Provider(), s.Kind(), s.Alarm.AlarmName)
}

// Provider returns the key reference to this provider
func (s *Alarm) Provider() string {
	return fmt.Sprintf("%s", Provider)
}

// Kind returns the key reference to this object
func (s *Alarm) Kind() string {
	return fmt.Sprintf("%s", AlarmKey)
}

// Execute a alarm raw string
func (s *Alarm) Execute(w io.Writer, tmpl *template.Template) error {
	return tmpl.Execute(w, s)
}

// Template returns the template
// func (s *Alarm) Template() string {
//     tmpl := `include {
//   path = "${find_in_parent_folders()}"
// }

// terraform {
//   source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-alarm"
// }

// inputs = {
//   # ESS Scaling Group
//   esssg_remote_state_bucket = "tkpd-tg-alicloud"
//   esssg_remote_state_key    = "{{ index .Extras "serviceName" }}/autoscale/ess-scaling-group/terraform.tfstate"

//   # ESS Scaling Rule
//   esssr_remote_state_bucket = "tkpd-tg-alicloud"
//   esssr_remote_state_key    = "{{ index .Extras "serviceName" }}/autoscale/ess-scaling-rules/{{ trimPrefix .ScalingRule.ScalingRuleName "tf-" }}/terraform.tfstate"

//   # ESS Alarm
//   essa_name                = "{{ trimPrefix .AlarmName "tf-" }}"
//   essa_enable              = {{ .Enable }}
//   essa_metric_type         = "{{ .MetricType }}"
//   essa_metric_name         = "{{ .MetricName }}"
//   essa_period              = {{ .Period }}
//   essa_statistics          = "{{ .Statistics }}"
//   essa_comparison_operator = "{{ .ComparisonOperator }}"
//   essa_threshold           = {{ .Threshold }}
//   essa_evaluation_count    = {{ .EvaluationCount }}
// }`

//     return tmpl
// }

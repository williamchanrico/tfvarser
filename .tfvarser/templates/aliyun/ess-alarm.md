# ess-alarm

## Available Fields

```
.Alarm
	AlarmName          string # e.g. "{{ .Alarm.AlarmName }}"
	AlarmID            string
	ScalingGroupID     string
	Enable             bool
	AlarmActions       []string
	MetricType         string
	MetricName         string
	Period             int
	Statistics         string
	ComparisonOperator string
	Threshold          float64
	EvaluationCount    int
	
.ScalingRule
	see ess-scaling-rule.md
	
.ScalingGroup
	see ess-scaling-group.md
	
.Extras
	map["serviceName"] # e.g. "{{ index .Extras "serviceName" }}"
```

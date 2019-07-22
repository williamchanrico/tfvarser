# ess-scaling-rule

## Available Fields

```
.ScalingRule
	ScalingRuleName string
	ScalingRuleID   string
	ScalingGroupID  string
	AdjustmentType  string
	AdjustmentValue int
	Cooldown        int

.ScalingGroup
	< see ess-scaling-group.md >
	
.Extras
	map["serviceName"] # e.g. "{{ index .Extras "serviceName" }}"
```

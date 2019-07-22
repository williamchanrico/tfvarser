# ess-lifecycle-hook

## Available Fields

```
.LifecycleHook
	LifecycleHookName   string
	LifecycleHookID     string
	LifecycleTransition string
	DefaultResult       string
	HeartbeatTimeout    int

.ScalingGroup
	< see ess-scaling-group.md >

.Extras
	map["serviceName"] # e.g. "{{ index .Extras "serviceName" }}"
```

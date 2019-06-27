<h1 align="center">tfvarser 👋</h1>
<p>
  <img src="https://img.shields.io/badge/version-0.1.0-blue.svg?cacheSeconds=2592000" />
</p>

> Generate tfvars file by mapping live cloud resources (query via SDK) to a tfvars template

A quick hack to reduce toils in importing hundreds of scaling groups by hand

*Disclaimer*: the templates are curated to specific needs and not really for general use, at least for now as this is only a quick hack to help current task.
So keeping the source close by will come in handy every now and then (modifying, rebuilding, etc.)

### Installation

`$ make go-build` will build the `tfvarser` binary in `./bin/tfvarser`

You can also build it yourself using simple `gotools`

### Usage & Examples

Requires some form of authentication to interact with cloud provider's API

#### Aliyun

```
export ALICLOUD_ACCESS_KEY=
export ALICLOUD_SECRET_KEY=
```

```
$ tfvarser -provider ali -obj ess -limit-names testapp

```

Command above will generate the following structure:
```
├── testapp
│   ├── autoscale
│   │   ├── ess-alarms
│   │   │   ├── go-testapp-downscale
│   │   │   │   └── terraform.tfvars
│   │   │   └── go-testapp-upscale
│   │   │       └── terraform.tfvars
│   │   ├── ess-lifecycle-hooks
│   │   │   ├── autoscaledown-event-mns-queue
│   │   │   │   └── terraform.tfvars
│   │   │   └── autoscaleup-event-mns-queue
│   │   │       └── terraform.tfvars
│   │   ├── ess-scaling-configurations
│   │   │   ├── go-testapp-1c-1gb
│   │   │   │   └── terraform.tfvars
│   │   │   └── go-testapp-1c-500mb
│   │   │       └── terraform.tfvars
│   │   ├── ess-scaling-group
│   │   │   └── terraform.tfvars
│   │   └── ess-scaling-rules
│   │       ├── auto-downscale
│   │       │   └── terraform.tfvars
│   │       └── auto-upscale
│   │           └── terraform.tfvars
```

Every provider objects like ESS or ECS in Aliyun may decide what the attribute that limit-names and limit-ids correspond to.

For example, in ESS Scaling Group, limit-names and limit-ids will limit by scaling group's name and ID as you might've guessed.

Currently, the tfvars templates are all hardcoded into the source.

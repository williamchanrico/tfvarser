package ess

import (
	"context"
	"fmt"
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// ScalingGroup struct to map into tfvars template
type ScalingGroup struct {
	ScalingGroupName string
	ScalingGroupID   string
	MinSize          int
	MaxSize          int
	RemovalPolicies  []string
	VSwitchIDs       []string
	MultiAZPolicy    string
}

// ScalingGroupTmpl is the tfvars template
var ScalingGroupTmpl = `
terragrunt = {
  include {
    path = "${find_in_parent_folders()}"
  }
  terraform {
    source = "git::git@github.com:tokopedia/tf-alicloud-modules.git//ess-scaling-group"
  }
}

# Name of the scaling group
esssg_name = "tf-{{ .ScalingGroupName }}"

# Minimum and maximum number of VMs in the scaling group
esssg_min_size = {{ .MinSize }}
esssg_max_size = {{ .MaxSize }}

# When downscaling, this specifies the order of VMs selected for removal
esssg_removal_policies = [
{{ range $index, $element := .RemovalPolicies }}{{- if $index }},
{{- end }}{{ if not $index }} {{ end }} "{{ $element -}}"{{ end }}
]

# VSwitches that will be used for created VMs, selection algorithm is based on esssg_multi_az_policy
esssg_vsw_ids          = [
{{ range $index, $element := .VSwitchIDs }}{{- if $index }},
{{- end }}{{ if not $index }} {{ end }} "{{ $element -}}"{{ end }}
]

# The order of VSwitches selected when creating new VMs
esssg_multi_az_policy  = "{{ .MultiAZPolicy }}"
`

// GetScalingGroups will query list of scaling groups
func (c *Client) GetScalingGroups(ctx context.Context) ([]ScalingGroup, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	responseChan := make(chan *esssdk.DescribeScalingGroupsResponse, 1)
	defer close(responseChan)

	var pageSize = 50

	// Get first page to calculate total number of pages to iterate
	c.getScalingGroupsByPage(ctx, responseChan, 1, pageSize)
	totalPageCount := ((<-responseChan).TotalCount / pageSize) + 1

	var wg sync.WaitGroup
	for pageNumber := totalPageCount; pageNumber > 0; pageNumber-- {
		wg.Add(1)
		go c.getScalingGroupsByPage(ctx, responseChan, pageNumber, pageSize)
	}

	scalingGroups := []ScalingGroup{}
	go func() {
		for resp := range responseChan {
			for _, sg := range resp.ScalingGroups.ScalingGroup {
				scalingGroup := ScalingGroup{}
				scalingGroup.ScalingGroupName = sg.ScalingGroupName
				scalingGroup.ScalingGroupID = sg.ScalingGroupId
				scalingGroup.MinSize = sg.MinSize
				scalingGroup.MaxSize = sg.MaxSize
				scalingGroup.RemovalPolicies = sg.RemovalPolicies.RemovalPolicy
				scalingGroup.VSwitchIDs = sg.VSwitchIds.VSwitchId
				scalingGroup.MultiAZPolicy = sg.MultiAZPolicy
				scalingGroups = append(scalingGroups, scalingGroup)
			}

			fmt.Printf("Retrived scaling groups, page number: %v\n", resp.PageNumber)
			wg.Done()
		}
	}()

	wg.Wait()

	return scalingGroups, nil
}

func (c *Client) getScalingGroupsByPage(ctx context.Context, responseChan chan *esssdk.DescribeScalingGroupsResponse, pageNumber int, pageSize int) {
	request := esssdk.CreateDescribeScalingGroupsRequest()
	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)

	respChan, errChan := c.ess.DescribeScalingGroupsWithChan(request)

	select {
	case resp := <-respChan:
		responseChan <- resp

	case _ = <-errChan:
		return

	case <-ctx.Done():
		return
	}
}

// GetScalingGroupInstances will return list of instances in a scaling group
func (c *Client) GetScalingGroupInstances(scalingGroupID string) ([]esssdk.ScalingInstance, error) {
	req := esssdk.CreateDescribeScalingInstancesRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var scalingInstanceList []esssdk.ScalingInstance

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.ess.DescribeScalingInstances(req)
		if err != nil {
			return nil, err
		}

		scalingInstanceList = append(scalingInstanceList, resp.ScalingInstances.ScalingInstance...)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.ScalingInstances.ScalingInstance))
	}

	return scalingInstanceList, nil
}

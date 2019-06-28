package ess

import (
	"context"
	"fmt"
	"time"

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

// GetScalingGroupsWithAsync will query list of scaling groups
func (c *Client) GetScalingGroupsWithAsync(ctx context.Context) ([]ScalingGroup, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var pageSize = 50
	var errCh = make(chan error, 1)

	// Get first page to calculate total number of pages to iterate
	firstPageCh := make(chan *esssdk.DescribeScalingGroupsResponse, 1)
	c.getScalingGroupsByPage(ctx, firstPageCh, errCh, 1, pageSize)
	select {
	case err := <-errCh:
		return nil, err
	default:
	}
	totalPageCount := ((<-firstPageCh).TotalCount / pageSize) + 1
	close(firstPageCh)

	// Scatter
	respCh := make(chan *esssdk.DescribeScalingGroupsResponse, totalPageCount)
	for pageNumber := 1; pageNumber <= totalPageCount; pageNumber++ {
		go c.getScalingGroupsByPage(ctx, respCh, errCh, pageNumber, pageSize)
	}

	// Gatter
	scalingGroups := []ScalingGroup{}
	for a := 0; a < totalPageCount; a++ {
		select {
		case resp := <-respCh:
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

		case err := <-errCh:
			return nil, err

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return scalingGroups, nil
}

func (c *Client) getScalingGroupsByPage(ctx context.Context, responseChan chan *esssdk.DescribeScalingGroupsResponse, errorChan chan error, pageNumber int, pageSize int) {
	request := esssdk.CreateDescribeScalingGroupsRequest()
	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)

	respChan, errChan := c.DescribeScalingGroupsWithChan(request)

	select {
	case resp := <-respChan:
		responseChan <- resp

	case err := <-errChan:
		errorChan <- err

	case <-ctx.Done():
		errorChan <- ctx.Err()
	}
}

// GetScalingGroupInstances will return list of instances in a scaling group
func (c *Client) GetScalingGroupInstances(scalingGroupID string) ([]esssdk.ScalingInstance, error) {
	req := esssdk.CreateDescribeScalingInstancesRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var scalingInstanceList []esssdk.ScalingInstance

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.DescribeScalingInstances(req)
		if err != nil {
			return nil, err
		}

		scalingInstanceList = append(scalingInstanceList, resp.ScalingInstances.ScalingInstance...)
		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.ScalingInstances.ScalingInstance))
	}

	return scalingInstanceList, nil
}

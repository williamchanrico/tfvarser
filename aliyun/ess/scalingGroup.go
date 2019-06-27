package ess

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

// GetScalingGroups will query list of scaling groups
func (c *Client) GetScalingGroups(ctx context.Context) ([]ScalingGroup, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	responseChan := make(chan *esssdk.DescribeScalingGroupsResponse, 1)
	defer close(responseChan)

	var pageSize = 50

	// Get first page to calculate total number of pages to iterate
	c.getScalingGroupsByPage(ctx, responseChan, 1, pageSize)
	totalPageCount := ((<-responseChan).TotalCount / pageSize) + 1
	if totalPageCount <= 0 {
		return nil, errors.New("Failed to retrieve scaling groups, please try again")
	}

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

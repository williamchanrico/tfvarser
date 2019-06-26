package ess

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	esssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

// Alarm struct to map to the tfvars template
type Alarm struct {
	AlarmName          string
	ScalingGroupID     string
	ScalingGroupName   string
	MetricType         string
	MetricName         string
	Period             int
	Statistics         string
	ComparisonOperator string
	Threshold          float64
	EvaluationCount    int
}

// GetAlarms returns list of scaling rule for the given scaling group
// scalingGroupName is only used to fill the struct, not for the request
func (c *Client) GetAlarms(scalingGroupID, scalingGroupName string) ([]Alarm, error) {
	req := esssdk.CreateDescribeAlarmsRequest()
	req.PageSize = requests.NewInteger(50)
	req.ScalingGroupId = scalingGroupID

	var alarms []Alarm

	for totalCount := req.PageSize; totalCount == req.PageSize; {
		resp, err := c.ess.DescribeAlarms(req)
		if err != nil {
			return nil, err
		}

		for _, al := range resp.AlarmList.Alarm {
			alarm := Alarm{}
			alarm.AlarmName = al.Name
			alarm.ScalingGroupID = al.ScalingGroupId
			alarm.ScalingGroupName = scalingGroupName // Needed for template
			alarm.MetricType = al.MetricType
			alarm.MetricName = al.MetricName
			alarm.Period = al.Period
			alarm.Statistics = al.Statistics
			alarm.ComparisonOperator = al.ComparisonOperator
			alarm.Threshold = al.Threshold
			alarm.EvaluationCount = al.EvaluationCount

			alarms = append(alarms, alarm)
		}

		req.PageNumber = requests.NewInteger(resp.PageNumber + 1)
		totalCount = requests.NewInteger(len(resp.AlarmList.Alarm))
	}

	return alarms, nil
}

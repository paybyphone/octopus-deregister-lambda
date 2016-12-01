package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Client struct {
	Ec2Client ec2iface.EC2API
}

func (c Client) GetEc2Instance(instanceId string) (*ec2.Instance, error) {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&instanceId},
	}
	response, err := c.Ec2Client.DescribeInstances(params)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve instance information for %s: %v", instanceId, err)
	}
	if len(response.Reservations) == 0 || len(response.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("No result found for %s", instanceId)
	}
	return response.Reservations[0].Instances[0], nil
}

func GetEc2InstanceTagValue(instance *ec2.Instance, tagName string) (string, bool) {
	for _, tag := range instance.Tags {
		if *tag.Key == tagName {
			return *tag.Value, true
		}
	}
	return "", false
}

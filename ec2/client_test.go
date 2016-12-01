package ec2_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/paybyphone/octopus-deregister-lambda/ec2"
)

type mockEc2Client struct {
	ec2iface.EC2API
}

func (m *mockEc2Client) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return mockOutput, mockError
}

var (
	mockOutput *ec2.DescribeInstancesOutput
	mockError  error
)

var _ = Describe("Ec2Client", func() {
	Describe("GetEc2Instance", func() {
		var (
			ec2Client      *mockEc2Client
			instanceResult *ec2.Instance
			errorResult    error
		)
		BeforeEach(func() {
			ec2Client = &mockEc2Client{}
		})
		JustBeforeEach(func() {
			client := Client{
				Ec2Client: ec2Client,
			}
			instanceResult, errorResult = client.GetEc2Instance("i-1234567890")
		})

		Context("When there is a result", func() {
			var instance ec2.Instance
			BeforeEach(func() {
				mockOutput = &ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{
						{
							Instances: []*ec2.Instance{
								&instance,
							},
						},
					},
				}
			})
			It("Returns the instance", func() {
				Expect(instanceResult).To(Equal(&instance))
			})
			It("Does not return an error", func() {
				Expect(errorResult).To(BeNil())
			})
		})

		Context("When there is no result", func() {
			Context("When no reservation is found", func() {
				BeforeEach(func() {
					mockOutput = &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{},
					}
				})
				It("Does not return an instance", func() {
					Expect(instanceResult).To(BeNil())
				})
				It("Returns an error", func() {
					Expect(errorResult).To(MatchError("No result found for i-1234567890"))
				})
			})

			Context("When no instance is found", func() {
				BeforeEach(func() {
					mockOutput = &ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{Instances: []*ec2.Instance{}},
						},
					}
				})
				It("Does not return an instance", func() {
					Expect(instanceResult).To(BeNil())
				})
				It("Returns an error", func() {
					Expect(errorResult).To(MatchError("No result found for i-1234567890"))
				})
			})
		})

		Context("When the AWS client returns an error", func() {
			BeforeEach(func() {
				mockOutput = nil
				mockError = fmt.Errorf("AWS failure!")
			})
			It("Does not return an instance", func() {
				Expect(instanceResult).To(BeNil())
			})
			It("Returns an error", func() {
				Expect(errorResult).To(MatchError("Failed to retrieve instance information for i-1234567890: AWS failure!"))
			})
		})
	})

	Describe("GetEc2InstanceTagValue", func() {
		Context("When the instance has no tags", func() {
			It("returns not found", func() {
				instance := ec2.Instance{}
				_, found := GetEc2InstanceTagValue(&instance, "the_tag")
				Expect(found).To(BeFalse())
			})
		})
		Context("When the instance has no matching tag", func() {
			It("returns not found", func() {
				key := "some_tag"
				value := "doesn't matter"
				instance := ec2.Instance{
					Tags: []*ec2.Tag{
						{Key: &key, Value: &value}},
				}
				_, found := GetEc2InstanceTagValue(&instance, "no_such_tag")
				Expect(found).To(BeFalse())
			})
		})
		Context("When the instance has a matching tag", func() {
			var (
				key           string
				expectedValue string
				actualValue   string
				actualFound   bool
			)
			BeforeEach(func() {
				key = "the_tag"
				expectedValue = "the value"
				instance := ec2.Instance{
					Tags: []*ec2.Tag{
						{Key: &key, Value: &expectedValue}},
				}
				actualValue, actualFound = GetEc2InstanceTagValue(&instance, key)
			})
			It("returns found", func() {
				Expect(actualFound).To(BeTrue())
			})
			It("returns the value", func() {
				Expect(actualValue).To(Equal(expectedValue))
			})
		})
	})
})

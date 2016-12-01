package main_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	octo "github.com/DimensionDataResearch/go-octo-api"

	"github.com/paybyphone/octopus-deregister-lambda"
)

type mockEc2Client struct {
	ec2iface.EC2API
}

func (m *mockEc2Client) DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return mockOutput, mockError
}

type mockOctopusClient struct{}

func (m mockOctopusClient) GetMachineByName(name string) (result *octo.Machine, found bool, err error) {
	return &octo.Machine{ID: "123"}, true, nil
}

func (m mockOctopusClient) DeleteMachine(machine *octo.Machine) (err error) {
	return nil
}

var (
	mockOutput *ec2.DescribeInstancesOutput
	mockError  error
)

var _ = Describe("Handler", func() {
	Describe("GetInstanceOctopusName", func() {
		var (
			handler     main.Handler
			actualName  string
			actualFound bool
			actualErr   error
		)
		BeforeEach(func() {
			ec2Client := &mockEc2Client{}
			handler = main.Handler{
				Ec2Client: ec2Client,
			}
		})
		JustBeforeEach(func() {
			actualName, actualFound, actualErr = handler.GetInstanceOctopusName("i-1234567890")
		})

		Context("When there is an instance", func() {
			Context("When there is a name", func() {
				expectedName := "the_octopus_name"
				BeforeEach(func() {
					key := "octopus_name"
					instance := ec2.Instance{
						Tags: []*ec2.Tag{
							{Key: &key, Value: &expectedName}},
					}
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
				It("Returns found", func() {
					Expect(actualFound).To(BeTrue())
				})
				It("Returns the name", func() {
					Expect(actualName).To(Equal(expectedName))
				})
				It("Does not return an error", func() {
					Expect(actualErr).ToNot(HaveOccurred())
				})
			})

			Context("When there is no name", func() {
				BeforeEach(func() {
					instance := ec2.Instance{}
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
				It("Returns not found", func() {
					Expect(actualFound).To(BeFalse())
				})
				It("Does not return an error", func() {
					Expect(actualErr).ToNot(HaveOccurred())
				})
			})
		})

		Context("When there is no instance", func() {
			BeforeEach(func() {
				mockOutput = &ec2.DescribeInstancesOutput{
					Reservations: []*ec2.Reservation{},
				}
			})
			It("Returns not found", func() {
				Expect(actualFound).To(BeFalse())
			})
			It("Does not return an error", func() {
				Expect(actualErr).To(MatchError("No result found for i-1234567890"))
			})
		})

		Context("When there is an AWS error", func() {
			BeforeEach(func() {
				mockError = fmt.Errorf("AWS failure!")
			})
			It("Returns not found", func() {
				Expect(actualFound).To(BeFalse())
			})
			It("Returns an error", func() {
				Expect(actualErr).To(MatchError("Failed to retrieve instance information for i-1234567890: AWS failure!"))
			})
		})
	})
})

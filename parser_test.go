package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/paybyphone/octopus-deregister-lambda"
)

var _ = Describe("Parser", func() {
	var eventJson []byte
	var eventEnvelope main.Ec2StateChangeNotificationEnvelope
	BeforeEach(func() {
		eventJson = []byte(`{
      "id": "7bf73129-1428-4cd3-a780-95db273d1602",
      "detail-type": "EC2 Instance State-change Notification",
      "source": "aws.ec2",
      "account": "123456789012",
      "time": "2015-11-11T21:29:54Z",
      "region": "us-east-1",
      "resources": [
        "arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111"
      ],
      "detail": {
        "instance-id": "i-abcd1111",
        "state": "pending"
      }
    }`)
	})
	JustBeforeEach(func() {
		var err error
		eventEnvelope, err = main.ParseEc2StateChangeNotification(eventJson)
		Expect(err).NotTo(HaveOccurred())
	})
	It("Parses the event envelope", func() {
		Expect(eventEnvelope.ID).To(Equal("7bf73129-1428-4cd3-a780-95db273d1602"))
		Expect(eventEnvelope.Region).To(Equal("us-east-1"))
	})
	It("Parses the event details", func() {
		Expect(eventEnvelope.Detail.InstanceID).To(Equal("i-abcd1111"))
		Expect(eventEnvelope.Detail.State).To(Equal("pending"))
	})
})

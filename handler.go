package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/sentientmonkey/future"

	octo "github.com/DimensionDataResearch/go-octo-api"
	pbp_ec2 "github.com/paybyphone/octopus-deregister-lambda/ec2"
	pbp_s3 "github.com/paybyphone/octopus-deregister-lambda/s3"
)

type Handler struct {
	Ec2Client           ec2iface.EC2API
	octopusClientFuture future.Future
}

type OctopusClient interface {
	GetMachineByName(name string) (result *octo.Machine, found bool, err error)
	DeleteMachine(machine *octo.Machine) (err error)
}

type Ec2StateChangeNotificationEnvelope struct {
	ID     string
	Region string
	Detail Ec2StateChangeNotificationDetail
}

type Ec2StateChangeNotificationDetail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

func ParseEc2StateChangeNotification(rawMessage json.RawMessage) (Ec2StateChangeNotificationEnvelope, error) {
	var eventEnvelope Ec2StateChangeNotificationEnvelope
	err := json.Unmarshal(rawMessage, &eventEnvelope)
	return eventEnvelope, err
}

func (h Handler) HandleEc2StateChange(eventDetail Ec2StateChangeNotificationDetail) (err error) {
	if eventDetail.State != "terminated" {
		log.Printf("Lambda misconfigured: received a \"%s\" event\n", eventDetail.State)
		// don't return an error, otherwise Lambda will retry
		return
	}

	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create AWS session: %v", err)
	}

	h.Ec2Client = ec2.New(sess)
	if err = h.StartConfiguringOctopusClient(sess); err != nil {
		return err
	}

	log.Printf("Getting Octopus name for instance \"%s\"", eventDetail.InstanceID)
	octopusName, found, err := h.GetInstanceOctopusName(eventDetail.InstanceID)
	if err != nil {
		return
	} else if !found {
		log.Printf("No Octopus name found for instance \"%s\"", eventDetail.InstanceID)
		return
	}
	log.Printf("Found Octopus name \"%s\" for instance \"%s\"", octopusName, eventDetail.InstanceID)

	octopusClient, err := h.OctopusClient()
	if err != nil {
		return
	}
	octopusMachine, found, err := octopusClient.GetMachineByName(octopusName)
	if err != nil {
		return
	} else if !found {
		log.Printf("No machine found for Octopus name \"%s\"", octopusName)
	}
	log.Printf("Found machine id \"%s\" for Octopus name \"%s\"", octopusMachine.ID, octopusName)

	err = octopusClient.DeleteMachine(octopusMachine)
	if err != nil {
		return
	}
	log.Printf("Deleted machine id \"%s\" from Octopus server", octopusMachine.ID)

	return
}

func (h Handler) GetInstanceOctopusName(instanceID string) (octopusName string, found bool, err error) {
	ec2Client := pbp_ec2.Client{
		Ec2Client: h.Ec2Client,
	}
	if os.Getenv("DEBUG_LOG") != "" {
		log.Printf("Getting EC2 Instance info for instance id \"%s\"", instanceID)
	}
	instance, err := ec2Client.GetEc2Instance(instanceID)
	if err != nil {
		log.Printf("Failed to get EC2 Instance info for instance id \"%s\": %v", instanceID, err)
		return
	}
	if os.Getenv("DEBUG_LOG") != "" {
		log.Printf("Got EC2 Instance info")
	}

	octopusName, found = pbp_ec2.GetEc2InstanceTagValue(instance, "octopus_name")

	return
}

func (h *Handler) StartConfiguringOctopusClient(sess *session.Session) error {
	s3cretsBucketName := os.Getenv("OCTOPUS_API_KEY_BUCKET")
	if s3cretsBucketName == "" {
		return errors.New("No OCTOPUS_API_KEY_BUCKET defined")
	}
	octopusApiKeyPath := os.Getenv("OCTOPUS_API_KEY_PATH")
	if octopusApiKeyPath == "" {
		return errors.New("No OCTOPUS_API_KEY_PATH defined")
	}
	s3Client := pbp_s3.Client{
		S3Client: s3.New(sess),
	}
	h.octopusClientFuture = future.NewFuture(func() (future.Value, error) {
		return h.CreateOctopusClient(s3Client, s3cretsBucketName, octopusApiKeyPath)
	})

	return nil
}

func (h Handler) OctopusClient() (OctopusClient, error) {
	client, err := h.octopusClientFuture.Get()
	if err != nil {
		return nil, err
	}
	return client.(OctopusClient), nil
}

func (h Handler) CreateOctopusClient(s3Client pbp_s3.Client, bucketName string, key string) (OctopusClient, error) {
	log.Println("Getting Octopus API key from S3 bucket…")
	octopusApiKey, err := s3Client.GetS3Object(bucketName, key)
	if err != nil {
		return nil, err
	}
	log.Println("Got Octopus API key.")

	serverUri := os.Getenv("OCTOPUS_URI")
	if serverUri == "" {
		return nil, errors.New("No OCTOPUS_URI defined")
	}

	return octo.NewClientWithAPIKey(serverUri, octopusApiKey)
}

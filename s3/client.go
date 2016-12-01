package s3

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Client struct {
	S3Client s3iface.S3API
}

func (c Client) GetS3ObjectReader(bucketName, key string) (io.ReadCloser, error) {
	params := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}
	response, err := c.S3Client.GetObject(params)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve S3 object %s/%s: %v", bucketName, key, err)
	}
	return response.Body, nil
}

func (c Client) GetS3Object(bucketName, key string) (string, error) {
	body, err := c.GetS3ObjectReader(bucketName, key)
	defer func() {
		if body != nil {
			body.Close()
		}
	}()
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("failed to read S3 object %s/%s: %v", bucketName, key, err)
	}
	return string(data[:]), nil
}

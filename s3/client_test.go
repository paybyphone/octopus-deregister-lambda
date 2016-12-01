package s3_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/paybyphone/octopus-deregister-lambda/s3"
)

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{Body: mockReader}, mockError
}

var (
	mockReader io.ReadCloser
	mockError  error
)

var _ = Describe("S3Client", func() {
	Describe("GetS3ObjectReader", func() {
		var (
			s3Client    *mockS3Client
			reader      io.ReadCloser
			errorResult error
		)
		BeforeEach(func() {
			s3Client = &mockS3Client{}
		})
		JustBeforeEach(func() {
			client := Client{
				S3Client: s3Client,
			}
			reader, errorResult = client.GetS3ObjectReader("bucket", "key")
		})

		Context("When there is a result", func() {
			BeforeEach(func() {
				mockReader = ioutil.NopCloser(strings.NewReader("foobar"))
				mockError = nil
			})
			It("Returns a reader", func() {
				content, err := ioutil.ReadAll(reader)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(Equal("foobar"))
			})
			It("Does not return an error", func() {
				Expect(errorResult).To(BeNil())
			})
		})

		Context("When the AWS client returns an error", func() {
			BeforeEach(func() {
				mockReader = nil
				mockError = fmt.Errorf("AWS failure!")
			})
			It("Does not return a reader", func() {
				Expect(reader).To(BeNil())
			})
			It("Returns an error", func() {
				Expect(errorResult).To(MatchError("failed to retrieve S3 object bucket/key: AWS failure!"))
			})
		})
	})

	Describe("GetS3Object", func() {
		var (
			s3Client    *mockS3Client
			content     string
			errorResult error
		)
		BeforeEach(func() {
			s3Client = &mockS3Client{}
		})
		JustBeforeEach(func() {
			client := Client{
				S3Client: s3Client,
			}
			content, errorResult = client.GetS3Object("bucket", "key")
		})

		Context("When there is a result", func() {
			BeforeEach(func() {
				mockReader = ioutil.NopCloser(strings.NewReader("foobar"))
				mockError = nil
			})
			It("Returns the content", func() {
				Expect(content).To(Equal("foobar"))
			})
			It("Does not return an error", func() {
				Expect(errorResult).To(BeNil())
			})
		})

		Context("When the AWS client returns an error", func() {
			BeforeEach(func() {
				mockReader = nil
				mockError = fmt.Errorf("AWS failure!")
			})
			It("Does not return any content", func() {
				Expect(content).To(Equal(""))
			})
			It("Returns an error", func() {
				Expect(errorResult).To(MatchError("failed to retrieve S3 object bucket/key: AWS failure!"))
			})
		})
	})
})

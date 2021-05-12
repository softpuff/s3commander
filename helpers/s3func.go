package helpers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSConfig struct {
	Session *session.Session
	Region  *string
}

type AWSOptions func(*AWSConfig)

func WithRegion(region string) AWSOptions {
	return func(a *AWSConfig) {
		a.Region = &region
	}
}

func NewAWSConfig(opts ...AWSOptions) (c *AWSConfig) {
	region := "us-west-2"

	c.Region = &region
	c.Session = session.Must(session.NewSession(&aws.Config{
		Region: c.Region,
	}))

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *AWSConfig) ListS3() (buckets []string, err error) {
	svc := s3.New(c.Session)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, fmt.Errorf("Creating bucker err: %v\n", err)
	}

	for _, b := range result.Buckets {
		buckets = append(buckets, *b.Name)
	}

	return
}

func (c *AWSConfig) ListS3Objects(b string, prefix string) (object []string, err error) {
	svc := s3.New(c.Session)

	input := s3.ListObjectsV2Input{
		Bucket: &b,
	}
	if prefix != "" {
		input.Prefix = &prefix
	}
	result, err := svc.ListObjectsV2(&input)
	if err != nil {
		return
	}

	for _, obj := range result.Contents {
		object = append(object, *obj.Key)
	}
	return
}

func (c *AWSConfig) PrintS3File(bucket string, key string) error {
	dl := s3manager.NewDownloader(c.Session)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	var buf []byte
	awsBuf := aws.NewWriteAtBuffer(buf)
	if _, err := dl.Download(awsBuf, params); err != nil {
		return err
	}

	fmt.Println(string(awsBuf.Bytes()))
	return nil

}

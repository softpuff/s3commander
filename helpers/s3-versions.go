package helpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (c *AWSConfig) GetVersions(bucket, prefix string) ([]*s3.ObjectVersion, error) {
	svc := s3.New(c.Session)

	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	result, err := svc.ListObjectVersions(input)
	if err != nil {
		return nil, err
	}
	return result.Versions, nil
}

func GetObjectVersions(ver []*s3.ObjectVersion) (vers []string) {
	for _, v := range ver {
		vers = append(vers, v.String())
	}
	return
}

func GetObjectVersionsIDs(ver []*s3.ObjectVersion) (ids []string) {
	for _, v := range ver {
		ids = append(ids, *v.VersionId)
	}
	return
}

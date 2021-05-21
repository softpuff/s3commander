package helpers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (c *AWSConfig) CountS3ObjectLines(bucket, key, expression string) error {

	svc := s3.New(c.Session)
	if expression == "" {
		expression = "select count(*) from S3Object s"
	}

	input := &s3.SelectObjectContentInput{
		Bucket:         &bucket,
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		Expression:     &expression,
		Key:            &key,
		InputSerialization: &s3.InputSerialization{
			JSON: &s3.JSONInput{
				Type: aws.String("Lines"),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			JSON: &s3.JSONOutput{},
		},
	}
	out, err := svc.SelectObjectContent(input)

	if err != nil {
		return err
	}

	defer out.EventStream.Close()

	for evt := range out.EventStream.Events() {
		switch e := evt.(type) {
		case *s3.RecordsEvent:
			fmt.Println(string(e.Payload))

		}
	}

	return nil
}

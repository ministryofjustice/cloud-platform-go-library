package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AwsCredentials represents the AWS credentials used to connect to an AWS account.
type AwsCredentials struct {
	Session *session.Session
	Profile string
	Region  string
}

// NewAwsCredentials constructs and populates a new AwsCredentials object
func NewAwsCreds(region string) (*AwsCredentials, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}

	return &AwsCredentials{
		Session: sess,
		Region:  region,
	}, nil
}

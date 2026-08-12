package backend

import (
	"context"
	"errors"
	"fmt"

	"github.com/AirHelp/treasury/backend/s3"
	"github.com/AirHelp/treasury/backend/ssm"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	ssmName = "ssm"
	s3Name  = "s3"
)

// Options for backend
type Options struct {
	Region       string
	S3BucketName string
	Backend      string
	AWSConfig    aws.Config
}

// New returns client for specific backend - s3 or ssm
// by default we use SSM
// once S3 bucket name is specified and no backend chosen we use S3
func New(options Options) (API, error) {
	if options.Backend == "" {
		if options.S3BucketName != "" {
			options.Backend = s3Name
		} else {
			options.Backend = ssmName
		}
	}
	switch options.Backend {
	case s3Name:
		return s3.New(options.Region, options.S3BucketName)
	case ssmName:
		awsConfig := options.AWSConfig
		if !isConfigured(awsConfig) {
			var err error
			awsConfig, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(options.Region))
			if err != nil {
				return nil, errors.Join(
					fmt.Errorf("unable to load SDK config with region %s", options.Region),
					err,
				)
			}
		}
		return ssm.New(awsConfig)
	}
	return nil, errors.New("invalid backend")
}

// isConfigured tells whether the caller provided its own AWS configuration.
// If it did not, the ambient one (environment, shared config, instance role)
// is loaded, the same way the S3 backend does it.
func isConfigured(awsConfig aws.Config) bool {
	return awsConfig.Region != "" || awsConfig.Credentials != nil
}

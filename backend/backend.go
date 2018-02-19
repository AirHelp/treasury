package backend

import (
	"errors"

	"github.com/AirHelp/treasury/backend/s3"
	"github.com/AirHelp/treasury/backend/ssm"
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
}

// New returns client for specific backend - s3 or ssm
// by default we use SSM
// once S3 bucket name is specified and no backend chosen we use S3
func New(options Options) (BackendAPI, error) {
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
		return ssm.New(options.Region)
	}
	return nil, errors.New("Invalid backend")
}

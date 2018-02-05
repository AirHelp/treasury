package backend

import (
	"errors"

	"github.com/AirHelp/treasury/backend/s3"
	"github.com/AirHelp/treasury/backend/ssm"
)

const defaultBackend = "ssm"

// Options for backend
type Options struct {
	Region       string
	S3BucketName string
	Backend      string
}

// New returns client for specific backend like s3
func New(options Options) (BackendAPI, error) {
	if options.Backend == "" {
		if options.S3BucketName != "" {
			options.Backend = "s3"
		} else {
			options.Backend = defaultBackend
		}
	}
	switch options.Backend {
	case "s3":
		return s3.New(options.Region, options.S3BucketName)
	case "ssm":
		return ssm.New(options.Region)
	}
	return nil, errors.New("Invalid backend")
}

package backend

import (
	"errors"

	"github.com/AirHelp/treasury/backend/s3"
)

// Options for backend
type Options struct {
	Region       string
	S3BucketName string
}

// New returns client for specific backend like s3
func New(options Options) (BackendAPI, error) {
	if options.S3BucketName == "" {
		return nil, errors.New("S3 bucket name is missing")
	}
	return s3.New(options.Region, options.S3BucketName)
}

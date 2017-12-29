package client_test

import (
	"testing"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/client"
	"github.com/AirHelp/treasury/test"
)

const (
	templateTestSourceFile      = "test/resources/source.secret.tpl"
	templateTestDestinationFile = "test/resources/destination.secret"
)

func TestTemplate(t *testing.T) {
	dummyClientOptions := &client.Options{
		AwsClient: &aws.Client{
			S3Svc: &test.MockS3Client{},
		},
	}
	treasury, err := client.New("fake_s3_bucket", dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	if err := treasury.Template(templateTestSourceFile, templateTestDestinationFile); err != nil {
		t.Error("Could not generate secret file from template. Error: ", err.Error())
	}
}

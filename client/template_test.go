package client_test

import (
	"os"
	"testing"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/client"
	"github.com/AirHelp/treasury/test"
)

const (
	templateTestSourceFile           = "../test/resources/source.secret.tpl"
	templateTestDestinationFile      = "../test/output/destination.secret"
	templateTestDestinationParentDir = "../test/output"
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
	if err := treasury.Template(templateTestSourceFile, templateTestDestinationFile, 0); err != nil {
		t.Error("Could not generate secret file from template. Error: ", err.Error())
	}
	_, err = os.Stat(templateTestDestinationParentDir)
	if err != nil {
		t.Error("Destination directory does not exist. Error: ", err.Error())
	}
}

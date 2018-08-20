package client_test

import (
	"os"
	"testing"

	"github.com/AirHelp/treasury/client"
	test "github.com/AirHelp/treasury/test/backend"
)

const (
	templateTestSourceFile           = "../test/resources/source.secret.tpl"
	templateTestDestinationFile      = "../test/output/destination.secret"
	templateTestDestinationParentDir = "../test/output"
)

func TestTemplate(t *testing.T) {
	dummyClientOptions := &client.Options{
		Backend:      &test.MockBackendClient{},
		S3BucketName: "fake_s3_bucket",
	}
	treasury, err := client.New(dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	if err := treasury.Template(templateTestSourceFile, templateTestDestinationFile, 0, map[string]string{}); err != nil {
		t.Error("Could not generate secret file from template. Error: ", err.Error())
	}
	_, err = os.Stat(templateTestDestinationParentDir)
	if err != nil {
		t.Error("Destination directory does not exist. Error: ", err.Error())
	}
}

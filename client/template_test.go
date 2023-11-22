package client_test

import (
	"os"
	"testing"

	"github.com/AirHelp/treasury/client"
	test "github.com/AirHelp/treasury/test/backend"
)

const (
	templateTestSourceFile           = "../test/resources/source.existing_secret.tpl"
	templateTestSourceFile2          = "../test/resources/source.not_existing_secret.tpl"
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

	envMap := map[string]string{
		"Environment": "test",
		"Name":        "some_testing_template",
	}

	tests := []struct {
		file    string
		wantErr bool
	}{
		{
			file:    templateTestSourceFile,
			wantErr: false,
		},
		{
			file:    templateTestSourceFile2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			err := treasury.Template(tt.file, templateTestDestinationFile, 0, map[string]string{}, envMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Failed to use treasury template, error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

	_, err = os.Stat(templateTestDestinationParentDir)
	if err != nil {
		t.Error("Destination directory does not exist. Error: ", err.Error())
	}
}

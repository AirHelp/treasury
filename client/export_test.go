package client_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/client"
	"github.com/AirHelp/treasury/test"
	"github.com/apex/treasury/cmd"
)

func TestExport(t *testing.T) {
	dummyClientOptions := &client.Options{
		AwsClient: &aws.Client{
			S3Svc: &test.MockS3Client{},
		},
	}
	treasury, err := client.New("fake_s3_bucket", dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	scenarios := []struct {
		key             string
		responseStrings []string
	}{
		{
			key: filepath.Dir(test.Key1) + "/",
			responseStrings: formatExportString(map[string]string{
				test.Key1: test.KeyValueMap[test.Key1],
				test.Key2: test.KeyValueMap[test.Key2],
			}),
		},
		{
			key: test.Key1,
			responseStrings: formatExportString(map[string]string{
				test.Key1: test.KeyValueMap[test.Key1],
			}),
		},
	}
	for _, scenario := range scenarios {
		exportString, err := treasury.Export(scenario.key, cmd.ExportString)
		if err != nil {
			t.Error(err)
		}
		for _, exportValue := range scenario.responseStrings {
			if !strings.Contains(exportString, exportValue) {
				t.Errorf("Wrong export string returned:\n%s, \nshoudl contain:\n%s", exportString, exportValue)
			}
		}
	}
}

func formatExportString(keyValue map[string]string) []string {
	var exportStrings []string
	for key, value := range keyValue {
		valueToExport := fmt.Sprintf(cmd.ExportString, filepath.Base(key), value)
		exportStrings = append(exportStrings, valueToExport)
	}
	return exportStrings
}

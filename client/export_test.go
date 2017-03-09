package client_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/client"
	"github.com/AirHelp/treasury/test"
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
		key            string
		responseString string
	}{
		{
			key: filepath.Dir(test.Key1) + "/",
			responseString: formatExportString(map[string]string{
				test.Key1: test.KeyValueMap[test.Key1],
				test.Key2: test.KeyValueMap[test.Key2],
			}),
		},
		{
			key: test.Key1,
			responseString: formatExportString(map[string]string{
				test.Key1: test.KeyValueMap[test.Key1],
			}),
		},
	}
	for _, scenario := range scenarios {
		exportString, err := treasury.Export(test.Key1)
		if err != nil {
			t.Error(err)
		}
		if exportString != scenario.responseString {
			t.Errorf("Wrong export string returned, \nexpected:%s, \ngot:%s", scenario.responseString, exportString)
		}
	}
}

func formatExportString(keyValue map[string]string) string {
	var buffer bytes.Buffer
	for key, value := range keyValue {
		buffer.WriteString(fmt.Sprintf(client.ExportString, filepath.Base(key), value))
	}
	response := buffer.String()
	return strings.Trim(response[:len(response)-2], " ")
}

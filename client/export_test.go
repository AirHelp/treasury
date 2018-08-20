package client_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/AirHelp/treasury/client"
	"github.com/AirHelp/treasury/cmd"
	test "github.com/AirHelp/treasury/test/backend"
)

func TestExport(t *testing.T) {
	dummyClientOptions := &client.Options{
		Backend:      &test.MockBackendClient{},
		S3BucketName: "fake_s3_bucket",
	}
	treasury, err := client.New(dummyClientOptions)
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
		exportString, err := treasury.Export(scenario.key, cmd.ExportString, map[string]string{})
		if err != nil {
			t.Error(err)
		}
		for _, exportValue := range scenario.responseStrings {
			if !strings.Contains(exportString, exportValue) {
				t.Errorf("Wrong export string returned:\n%s, \nshould contain:\n%s", exportString, exportValue)
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

func TestClient_ExportMap(t *testing.T) {
	dummyClientOptions := &client.Options{
		Backend: &test.MockBackendClient{},
	}
	c, err := client.New(dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name       string
		key        string
		wantResult map[string]string
		wantErr    bool
	}{
		{
			name: "correct key",
			key:  "test/webapp/",
			wantResult: map[string]string{
				test.ShortKey1: test.KeyValueMap[test.Key1],
				test.ShortKey2: test.KeyValueMap[test.Key2],
			},
			wantErr: false,
		},
		{
			name: "key with only 1 result",
			key:  "test/cockpit/",
			wantResult: map[string]string{
				test.ShortKey3: test.KeyValueMap[test.Key3],
			},
			wantErr: false,
		},
		{
			name:       "wrong key to export",
			key:        "test/webapp",
			wantResult: map[string]string{},
			wantErr:    true,
		},
		{
			name:       "correct key with no results",
			key:        "test/dummyApplication/",
			wantResult: map[string]string{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := c.ExportMap(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ExportMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Client.ExportMap() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestClient_ExportToTemplate(t *testing.T) {
	dummyClientOptions := &client.Options{
		Backend: &test.MockBackendClient{},
	}
	c, err := client.New(dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name      string
		key       string
		want      string
		appendMap map[string]string
		wantErr   bool
	}{
		{
			name:    "full key path",
			key:     test.Key1,
			want:    fmt.Sprintf("%s=%s\n", test.ShortKey1, test.KeyValueMap[test.Key1]),
			wantErr: false,
		},
		{
			name: "correct prefix key",
			key:  "test/webapp/",
			want: fmt.Sprintf("%s=%s\n%s=%s\n",
				test.ShortKey1, test.KeyValueMap[test.Key1],
				test.ShortKey2, test.KeyValueMap[test.Key2],
			),
			wantErr: false,
		},
		{
			name:      "merged variable",
			key:       "test/airmail/",
			appendMap: map[string]string{"DATABASE_URL": "?pool=10"},
			want: fmt.Sprintf("%s=%s\n%s=%s\n",
				test.ShortKey4, test.KeyValueMap[test.Key4]+"?pool=10",
				test.ShortKey5, test.KeyValueMap[test.Key5],
			),
			wantErr: false,
		},
		{
			name:      "merged variable - multiple vars",
			key:       "test/aircom/",
			appendMap: map[string]string{"NEW_RELIC_LICENSE_KEY": "test2", "TWILIO_AUTH_TOKEN": "test1"},
			want: fmt.Sprintf("%s=%s\n%s=%s\n",
				test.ShortKey7, test.KeyValueMap[test.Key7]+"test2",
				test.ShortKey6, test.KeyValueMap[test.Key6]+"test1",
			),
			wantErr: false,
		},
		{
			name:    "incorrect prefix key",
			key:     "bla_bla_bla",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.ExportToTemplate(tt.key, tt.appendMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ExportToTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Client.ExportToTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

package client_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/AirHelp/treasury/aws"
	"github.com/AirHelp/treasury/client"
	"github.com/AirHelp/treasury/test"
)

func TestRead(t *testing.T) {
	dummyClientOptions := &client.Options{
		AwsClient: &aws.Client{
			S3Svc: &test.MockS3Client{},
		},
	}
	treasury, err := client.New("fake_s3_bucket", dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	secret, err := treasury.Read(test.Key1)
	if err != nil {
		t.Error(err)
	}
	if secret.Value != test.KeyValueMap[test.Key1] {
		t.Errorf("Read returns wrong secret")
	}
}

func TestReadValue(t *testing.T) {
	dummyClientOptions := &client.Options{
		AwsClient: &aws.Client{
			S3Svc: &test.MockS3Client{},
		},
	}
	treasury, err := client.New("fake_s3_bucket", dummyClientOptions)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "test valid key",
			key:     test.Key1,
			want:    test.KeyValueMap[test.Key1],
			wantErr: false,
		},
		{
			name:    "test non existing key",
			key:     "nonExistingKey",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := treasury.ReadValue(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ReadValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.ReadValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadGroup(t *testing.T) {
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
		responseSecrets int
	}{
		{
			key:             filepath.Dir(test.Key1) + "/",
			responseSecrets: 2,
		},
		{
			key:             test.Key1,
			responseSecrets: 1,
		},
	}
	for _, scenario := range scenarios {
		secrets, err := treasury.ReadGroup(scenario.key)
		if err != nil {
			t.Error(err)
		}
		if len(secrets) != scenario.responseSecrets {
			t.Errorf("Wrong number of returned secrets, expected:%d, got:%d", scenario.responseSecrets, len(secrets))
		}
		for _, secret := range secrets {
			if !strings.Contains(secret.Key, scenario.key) {
				t.Errorf("Secret key:%s should contains argument key:%s", secret.Key, scenario.key)
			}
			if secret.Value != test.KeyValueMap[secret.Key] {
				t.Errorf("Wrong value for key:%s, expected:%s, got:%s", secret.Key, test.KeyValueMap[secret.Key], secret.Value)
			}
		}
	}

}

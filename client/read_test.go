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
		t.Errorf("Reads returns wrong secret")
	}
}

func TestReadGroupd(t *testing.T) {
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

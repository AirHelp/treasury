package backend

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AirHelp/treasury/backend"
	"github.com/AirHelp/treasury/types"
)

const (
	Key1      = "test/webapp/cockpit_api_pass"
	Key1NoEnv = "webapp/cockpit_api_pass"
	ShortKey1 = "cockpit_api_pass"
	Key2      = "test/webapp/user_api_pass"
	Key2NoEnv = "webapp/user_api_pass"
	ShortKey2 = "user_api_pass"
	Key3      = "test/cockpit/user_api_pass"
	Key3NoEnv = "cockpit/user_api_pass"
	ShortKey3 = "user_api_pass"
	Key4      = "test/webapp/some_key"
	Key4NoEnv = "webapp/some_key"
	ShortKey4 = "some_key"
	Key5      = "test/airmail/DATABASE_URL"
	Key5NoEnv = "airmail/DATABASE_URL"
	ShortKey5 = "DATABASE_URL"
	Key6      = "test/airmail/user_api_pass"
	Key6NoEnv = "airmail/user_api_pass"
	ShortKey6 = "user_api_pass"
	Key7      = "test/aircom/TWILIO_AUTH_TOKEN"
	Key7NoEnv = "aircom/TWILIO_AUTH_TOKEN"
	ShortKey7 = "TWILIO_AUTH_TOKEN"
	Key8      = "test/aircom/NEW_RELIC_LICENSE_KEY"
	Key8NoEnv = "aircom/NEW_RELIC_LICENSE_KEY"
	ShortKey8 = "NEW_RELIC_LICENSE_KEY"
)

var KeyValueMap = map[string]string{
	Key1: "as9@#$%^&*(/2hdiwnf",
	Key2: "as9@#$&*(/2saddsahdiwnf",
	Key3: "#$&*(/2saddsah&as",
	Key4: "value=with=multiple=equal=signs==",
	Key5: "postgres://user:password@ip:port/db",
	Key6: "2oui3yrwohsf",
	Key7: "weoirgfhdh",
	Key8: "sfjsoidhgi340j",
}

// MockBackendClient fake backendAPI
type MockBackendClient struct {
	backend.BackendAPI
}

func (m *MockBackendClient) PutObject(input *types.PutObjectInput) error {
	if _, ok := KeyValueMap[input.Key]; !ok {
		return errors.New(fmt.Sprintf("Missing key:%s in KeyValue map", input.Key))
	}
	return nil
}

func (m *MockBackendClient) GetObject(input *types.GetObjectInput) (*types.GetObjectOutput, error) {
	return &types.GetObjectOutput{
		Value: KeyValueMap[input.Key],
	}, nil
}

func (m *MockBackendClient) GetObjects(input *types.GetObjectsInput) (*types.GetObjectsOuput, error) {
	response := make(map[string]string)
	for key := range KeyValueMap {
		if strings.Contains(key, input.Prefix) {
			response[key] = KeyValueMap[key]
		}
	}
	return &types.GetObjectsOuput{Secrets: response}, nil
}

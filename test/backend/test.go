package backend

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/AirHelp/treasury/backend"
	"github.com/AirHelp/treasury/types"
)

const (
	Key1 = "test/webapp/cocpit_api_pass"
	Key2 = "test/webapp/user_api_pass"
	Key3 = "test/cockpit/user_api_pass"
)

var KeyValueMap = map[string]string{
	Key1: "as9@#$%^&*(/2hdiwnf",
	Key2: "as9@#$&*(/2saddsahdiwnf",
	Key3: "#$&*(/2saddsah&as",
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
		Body: ioutil.NopCloser(
			bytes.NewReader(
				[]byte(KeyValueMap[input.Key]))),
	}, nil
}

func (m *MockBackendClient) GetObjects(input *types.GetObjectsInput) (*types.GetObjectsOuput, error) {
	response := make(map[string]string)
	for key := range KeyValueMap {
		if strings.Contains(key, input.Prefix) {
			response[key] = KeyValueMap[key]
		}
	}
	return &types.GetObjectsOuput{response}, nil
}

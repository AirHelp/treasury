package ssm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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

const (
	SSMKey1 = "/" + Key1
	SSMKey2 = "/" + Key2
	SSMKey3 = "/" + Key3
)

var SSMKeyValueMap = map[string]string{
	SSMKey1: KeyValueMap[Key1],
	SSMKey2: KeyValueMap[Key2],
	SSMKey3: KeyValueMap[Key3],
}

// MockSSMClient is a mock implementation of the SSMClient interface.
type MockSSMClient struct {
	Parameters map[string]string
}

type SSMClient interface {
	PutParameter(ctx context.Context, input *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
	GetParameter(ctx context.Context, input *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	GetParametersByPath(ctx context.Context, input *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
	DeleteParameter(ctx context.Context, input *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error)
}

//PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
func (m *MockSSMClient) PutParameter(ctx context.Context, input *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error) {
	if input == nil {
		return nil, fmt.Errorf("PutParameterInput is empty")
	}
	if *input.Name == "" {
		return nil, fmt.Errorf("Name in PutParameterInput is not set")
	}
	// SSM path based key needs to start from "/"
	name := *input.Name
	if string(name[0]) != "/" {
		return nil, fmt.Errorf("SSM Name needs to start from /")
	}
	if *input.Value == "" {
		return nil, fmt.Errorf("Value in PutParameterInput is not set")
	}
	var version int64 = 1
	return &ssm.PutParameterOutput{Version: version}, nil
}

//GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
func (m *MockSSMClient) GetParameter(ctx context.Context, input *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	log.Println("input.Name:", input.Name)
	if _, ok := SSMKeyValueMap[*input.Name]; !ok {
		return nil, fmt.Errorf("Missing key:%s in KeyValue map", *input.Name)
	}
	if !*input.WithDecryption {
		return nil, fmt.Errorf("Missing decryption field")
	}
	value := SSMKeyValueMap[*input.Name]
	return &ssm.GetParameterOutput{
		// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#Parameter
		Parameter : &types.Parameter{
			Name:  input.Name,
			Value: &value,
		},
	}, nil
}

// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#SSM.GetParametersByPath
// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#GetParametersByPathInput
func (m *MockSSMClient) GetParametersByPath(ctx context.Context, input *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error) {
	if !*input.WithDecryption {
		return nil, fmt.Errorf("Missing decryption field")
	}
	var contents []types.Parameter
	for key, value := range SSMKeyValueMap {
		key := key
		value := value
		if strings.Contains(key, *input.Path) {
			contents = append(contents, types.Parameter{
				Name:  &key,
				Value: &value,
			})
		}
	}
	return &ssm.GetParametersByPathOutput{
		Parameters: contents,
	}, nil
}

func (m *MockSSMClient) DeleteParameter(ctx context.Context, input *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error) {
	if input == nil {
		return nil, fmt.Errorf("DeleteParameterInput is empty")
	}
	if *input.Name == "" {
		return nil, fmt.Errorf("Name in DeleteParameterInput is not set")
	}
	// SSM path based key needs to start from "/"
	name := *input.Name
	if string(name[0]) != "/" {
		return nil, fmt.Errorf("SSM Name needs to start from /")
	}
	return &ssm.DeleteParameterOutput{}, nil
}
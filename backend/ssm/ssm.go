package ssm

import (
	"context"
	"errors"

	"github.com/AirHelp/treasury/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

const defaultParameterType = "SecureString"

// PutObject writes a given secret value on SSM
// it uses PutParameter API call
// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_PutParameter.html
func (c *Client) PutObject(object *types.PutObjectInput) error {
	if object.Key == "" {
		return errors.New("The key name is not valid.")
	}
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ssm#PutParameterInput
	putParameterInput := &ssm.PutParameterInput{
		KeyId: aws.String("alias/" + object.Environment),
		// we decided to use path based keys without `/` at the beginning
		// so we need to add it here
		Name:      aws.String("/" + object.Key),
		Type:      ssmtypes.ParameterType(defaultParameterType),
		Value:     aws.String(object.Value),
		Overwrite: aws.Bool(true),
	}

	// PutParameter returns Version of the parameter
	// shall we validate this version?
	_, err := c.svc.PutParameter(context.Background(), putParameterInput)
	return err
}

// GetObject returns a secret for given key
func (c *Client) GetObject(object *types.GetObjectInput) (*types.GetObjectOutput, error) {
	params := &ssm.GetParameterInput{
		// we decided to use path based keys without `/` at the beginning
		// so we need to add it here
		Name: aws.String("/" + object.Key),
		// Retrieve all parameters in a hierarchy with their value decrypted.
		WithDecryption: aws.Bool(true),
	}

	resp, err := c.svc.GetParameter(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return &types.GetObjectOutput{Value: *resp.Parameter.Value}, nil
}

// GetObjects returns key value map for given pattern/prefix
func (c *Client) GetObjects(object *types.GetObjectsInput) (*types.GetObjectsOuput, error) {
	var nextToken *string
	var parameters []ssmtypes.Parameter
	for {
		// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ssm#GetParametersByPathInput
		getParametersByPathInput := &ssm.GetParametersByPathInput{
			Path: aws.String("/" + object.Prefix),
			// Retrieve all parameters in a hierarchy with their value decrypted.
			WithDecryption: aws.Bool(true),
			MaxResults:     aws.Int32(10),
			NextToken:      nextToken,
		}

		// we're only interested with GetParametersByPathOutput.Parameters
		// Parameters []Parameter `type:"list"`
		// See also, https://docs.aws.amazon.com/goto/WebAPI/ssm-2014-11-06/Parameter
		resp, err := c.svc.GetParametersByPath(context.Background(), getParametersByPathInput)
		if err != nil {
			return nil, err
		}

		for _, parameter := range resp.Parameters {
			parameters = append(parameters, parameter)
		}
		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	keyValuePairs := make(map[string]string, len(parameters))
	for _, parameter := range parameters {
		keyValuePairs[unSlash(*parameter.Name)] = *parameter.Value
	}
	return &types.GetObjectsOuput{Secrets: keyValuePairs}, nil
}

// unSlash removes 1st char from a string
// GetParametersByPath from SSM returns key path with "/" at the beginning
// but we don't need it :)
func unSlash(input string) string {
	if string(input[0]) == "/" {
		return input[1:]
	}
	return input
}

func (c *Client) DeleteObject(object *types.DeleteObjectInput) error {
	params := &ssm.DeleteParameterInput{
		// we decided to use path based keys without `/` at the beginning
		// so we need to add it here
		Name: aws.String("/" + object.Key),
	}
	_, err := c.svc.DeleteParameter(context.Background(), params)
	return err
}

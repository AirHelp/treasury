package ssm

import (
	"github.com/AirHelp/treasury/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const defaultParameterType = "SecureString"

// PutObject writes a given secret value on SSM
// it uses PutParameter API call
// https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_PutParameter.html
func (c *Client) PutObject(object *types.PutObjectInput) error {
	// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#PutParameterInput
	putParameterInput := &ssm.PutParameterInput{
		KeyId: aws.String("alias/" + object.Environment),
		// we decided to use path based keys without `/` at the begining
		// so we need to add it here
		Name:      aws.String("/" + object.Key),
		Type:      aws.String(defaultParameterType),
		Value:     aws.String(object.Value),
		Overwrite: aws.Bool(true),
	}

	// PutParameter returns Version of the parameter
	// shall we validate this version?
	_, err := c.svc.PutParameter(putParameterInput)
	if err != nil {
		return err
	}

	return nil
}

// GetObject returns a secret for given key
func (c *Client) GetObject(object *types.GetObjectInput) (*types.GetObjectOutput, error) {
	params := &ssm.GetParameterInput{
		// we decided to use path based keys without `/` at the begining
		// so we need to add it here
		Name: aws.String("/" + object.Key),
		// Retrieve all parameters in a hierarchy with their value decrypted.
		WithDecryption: aws.Bool(true),
	}

	resp, err := c.svc.GetParameter(params)
	if err != nil {
		return nil, err
	}

	return &types.GetObjectOutput{Value: *resp.Parameter.Value}, nil
}

// GetObjects returns key value map for given pattern/prefix
func (c *Client) GetObjects(object *types.GetObjectsInput) (*types.GetObjectsOuput, error) {
	// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#SSM.GetParametersByPath
	getParametersByPathInput := &ssm.GetParametersByPathInput{
		Path: aws.String("/" + object.Prefix),
		// Retrieve all parameters in a hierarchy with their value decrypted.
		WithDecryption: aws.Bool(true),
	}

	// we're only interested with GetParametersByPathOutput.Parameters
	// Parameters []*Parameter `type:"list"`
	// See also, https://docs.aws.amazon.com/goto/WebAPI/ssm-2014-11-06/Parameter
	resp, err := c.svc.GetParametersByPath(getParametersByPathInput)
	if err != nil {
		return nil, err
	}

	keyValuePairs := make(map[string]string, len(resp.Parameters))
	for _, parameter := range resp.Parameters {
		keyValuePairs[unSlash(*parameter.Name)] = *parameter.Value
	}
	return &types.GetObjectsOuput{Secrets: keyValuePairs}, nil
}

// unShash removes 1st char from a string
// GetParametersByPath from SSM returns key path with "/" at the beginning
// but we don't need it :)
func unSlash(input string) string {
	if string(input[0]) == "/" {
		return input[1:]
	}
	return input
}

package aws

import "github.com/aws/aws-sdk-go/service/iam"

// GetUserName returns user name based on the AWS access key ID
func (c *Client) GetUserName() (string, error) {
	svc := iam.New(c.sess)
	params := &iam.GetUserInput{}
	resp, err := svc.GetUser(params)
	if err != nil {
		return "", err
	}

	return *resp.User.UserName, nil
}

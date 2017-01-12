package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// Encrypt with KMS
// http://docs.aws.amazon.com/kms/latest/APIReference/API_Encrypt.html
func (c *Client) Encrypt(kmsAlias, secret string) (*kms.EncryptOutput, error) {
	svc := kms.New(c.sess)
	params := &kms.EncryptInput{
		KeyId:     aws.String("alias/" + kmsAlias),
		Plaintext: []byte(secret),
		// TO DO: add context support
		// EncryptionContext: map[string]*string{
		// 	"Key": aws.String("EncryptionContextValue"), // Required
		// 	// More values...
		// },
	}
	return svc.Encrypt(params)
}

// Decrypt with KMS
// http://docs.aws.amazon.com/kms/latest/APIReference/API_Decrypt.html
func (c *Client) Decrypt(blobSecret []byte) (string, error) {
	svc := kms.New(c.sess)

	params := &kms.DecryptInput{
		CiphertextBlob: []byte(blobSecret),
		// TO DO: add support for context
		// EncryptionContext: map[string]*string{
		//     "Key": aws.String("EncryptionContextValue"),
		//     // More values...
		// },
	}
	resp, err := svc.Decrypt(params)
	if err != nil {
		return "", err
	}

	return string(resp.Plaintext), nil
}

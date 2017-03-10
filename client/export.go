package client

import (
	"bytes"
	"fmt"
	"path/filepath"
)

// Export returns command exporting found secrets
func (c *Client) Export(key, singleKeyExportFormat string) (string, error) {
	secrets, err := c.ReadGroup(key)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	for _, secret := range secrets {
		buffer.WriteString(fmt.Sprintf(singleKeyExportFormat, filepath.Base(secret.Key), secret.Value))
	}
	return buffer.String(), nil
}

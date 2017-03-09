package client

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	ExportString = " export %s='%s' &&"
)

// Export returns decrypted secret for given key
func (c *Client) Export(key string) (string, error) {
	secrets, err := c.ReadGroup(key)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	for _, secret := range secrets {
		buffer.WriteString(fmt.Sprintf(ExportString, filepath.Base(secret.Key), secret.Value))
	}
	exportCommand := strings.Trim(buffer.String()[:len(buffer.String())-2], " ")
	return exportCommand, nil
}

package client

import (
	"bytes"
	"fmt"
	"path/filepath"
)

const (
	// ExportString format of single export string
	ExportString = "export %s='%s'\n"
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
	exportCommand := buffer.String()[:len(buffer.String())]
	return exportCommand, nil
}

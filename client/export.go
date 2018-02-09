package client

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
)

// Export returns secrets in given format
// format should be provided in singleKeyExportFormat
// e.g.: singleKeyExportFormat = "export %s='%s'\n"
func (c *Client) Export(key, singleKeyExportFormat string) (string, error) {
	secrets, err := c.ReadGroup(key)
	if err != nil {
		return "", err
	}
	var sortedKeys []string
	keySecretMap := make(map[string]*Secret, len(secrets))
	for _, secret := range secrets {
		sortedKeys = append(sortedKeys, secret.Key)
		keySecretMap[secret.Key] = secret
	}
	sort.Strings(sortedKeys)
	var buffer bytes.Buffer
	for _, key := range sortedKeys {
		secret := keySecretMap[key]
		buffer.WriteString(fmt.Sprintf(singleKeyExportFormat, filepath.Base(secret.Key), secret.Value))
	}
	return buffer.String(), nil
}

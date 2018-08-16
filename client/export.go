package client

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AirHelp/treasury/utils"
)

// Export returns secrets in given format
// format should be provided in singleKeyExportFormat
// e.g.: singleKeyExportFormat = "export %s='%s'\n"
func (c *Client) Export(key, singleKeyExportFormat string) (string, error) {
	var secrets []*Secret
	var err error
	// if we get valid prefix we use ReadGroup method
	// else we have 1 key only and we use Read method
	if validPrefix(key) {
		secrets, err = c.ReadGroup(key)
		if err != nil {
			return "", err
		}
	} else {
		secret, err := c.Read(key)
		if err != nil {
			return "", err
		}
		secrets = append(secrets, secret)
	}
	var sortedKeys []string
	keySecretMap := make(map[string]*Secret, len(secrets))
	for _, secret := range secrets {
		sortedKeys = append(sortedKeys, secret.Key)
		keySecretMap[secret.Key] = secret
	}

	var AppendMap map[string]string
	AppendMap = make(map[string]string)
	for _, val := range c.Append {
		parts := strings.Split(val, ":")
		if len(parts) == 2 {
			AppendMap[parts[0]] = parts[1]
		} else {
			return "", errors.New("Bad append format (--append <variable>:<string>)")
		}
	}

	sort.Strings(sortedKeys)
	var buffer bytes.Buffer
	for _, key := range sortedKeys {
		secret := keySecretMap[key]
		secret.Value = fmt.Sprintf("%s%s", secret.Value, AppendMap[filepath.Base(secret.Key)])
		buffer.WriteString(fmt.Sprintf(singleKeyExportFormat, filepath.Base(secret.Key), secret.Value))
	}
	return buffer.String(), nil
}

func (c *Client) ExportToTemplate(key string) (string, error) {
	return c.Export(key, "%s=%s\n")
}

// validPrefix returns true if correct Prefix is given as an input
// e.g.: test/key/ is an correct prefix
// test/key/var is a full key path not a prefix
func validPrefix(input string) bool {
	err := utils.ValidateInputKey(input)
	return (err != nil) == true
}

// ExportMap returns map of Key=Value secrets (Key is without full path)
func (c *Client) ExportMap(key string) (map[string]string, error) {
	results := make(map[string]string)
	secrets, err := c.ReadGroup(key)
	if err != nil {
		return results, err
	}
	for _, secret := range secrets {
		splitKey := strings.Split(secret.Key, "/")
		results[splitKey[2]] = secret.Value
	}
	return results, nil
}

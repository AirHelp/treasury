package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const (
	// DefaultDestinationFilePerms are the default file permissions for destination file rendered into
	// disk when a specific file permission has not already been specified.
	DefaultDestinationFilePerms = 0644
)

// Read template file
func readTemplate(filePath string) (string, error) {
	templateBuffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(templateBuffer), nil
}

// Render template
func (c *Client) renderTemplate(templateText string, appendMap, envMap map[string]string) (templateResultBuffer bytes.Buffer, err error) {
	// Create a FuncMap with which to register the function.
	funcMap := template.FuncMap{
		"readFromEnv": c.ReadFromEnv,
		"exportMap":   c.ExportMap,
		// The name "read" is what the function will be called in the template text.
		"read": func(key string) (string, error) {
			return c.ReadValue(key)
		},
		"export": func(key string) (string, error) {
			return c.ExportToTemplate(key, appendMap)
		},
		"exportFromEnv": func(environment, key string) (string, error) {
			return c.ExportToTemplate(fmt.Sprintf("%s/%s/", environment, key), appendMap)
		},
	}
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("templateCli").Funcs(funcMap).Parse(templateText)
	if err != nil {
		return
	}
	// Run the template.
	err = tmpl.Execute(&templateResultBuffer, envMap)
	return
}

// Write template results to the destination file
func writeTemplateResults(destinationFilePath string, templateResultBuffer bytes.Buffer, perms os.FileMode) error {
	// If the parent destination directory does not exist, it will be created
	// automatically with permissions 0755. To use a different permission, create
	// the directory first or use `chmod` in a Command.
	parent := filepath.Dir(destinationFilePath)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		if err := os.MkdirAll(parent, 0755); err != nil {
			return err
		}
	}

	// If the user did not explicitly set permissions, attempt to lookup the
	// current permissions on the file. If the file does not exist, fall back to
	// the default. Otherwise, inherit the current permissions.
	if perms == 0 {
		currentDestinationFileInfo, err := os.Stat(destinationFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				perms = DefaultDestinationFilePerms
			} else {
				return err
			}
		} else {
			perms = currentDestinationFileInfo.Mode()
		}
	}

	return ioutil.WriteFile(destinationFilePath, templateResultBuffer.Bytes(), perms)
}

// Template generates a file with secrets from given template
func (c *Client) Template(sourceFilePath, destinationFilePath string, perms os.FileMode, appendMap, envMap map[string]string) error {
	templateText, err := readTemplate(sourceFilePath)
	if err != nil {
		return err
	}
	templateResultBuffer, err := c.renderTemplate(templateText, appendMap, envMap)
	if err != nil {
		return err
	}
	return writeTemplateResults(destinationFilePath, templateResultBuffer, perms)
}

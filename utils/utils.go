package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	environmentApplicationKeyRegex = `^[a-zA-Z-_]+\/[a-zA-Z-_0-9]+\/[a-zA-Z-_0-9]+$`
	environmentApplicationRegex    = `^[a-zA-Z-_]+\/[a-zA-Z-_0-9]+\/[a-zA-Z-_0-9]*$`
	userUsernameKeyRegex           = `^user\/[a-zA-Z]+\.[a-zA-Z]+\/[a-zA-Z-_0-9]+$`
	userUsernameRegex              = `^user\/[a-zA-Z]+\.[a-zA-Z]+\/[a-zA-Z-_0-9]*$`
)

// ValidateInputKey checks if cli input is valid
func ValidateInputKey(cliIn string) error {
	if strings.HasPrefix(cliIn, "user") {
		return validate(cliIn, userUsernameKeyRegex)
	}
	return validate(cliIn, environmentApplicationKeyRegex)
}

// ValidateInputKeyPattern checks if cli input is valid, also without key name
func ValidateInputKeyPattern(cliIn string) error {
	if strings.HasPrefix(cliIn, "user") {
		return validate(cliIn, userUsernameRegex)
	}
	return validate(cliIn, environmentApplicationRegex)
}

func validate(cliIn, pattern string) error {
	match, err := regexp.MatchString(pattern, cliIn)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("Given Key (%s) does not match to defined regex.", cliIn)
	}
	return nil
}

// FindEnvironmentApplicationName slices cli input into environment and application names
func FindEnvironmentApplicationName(cliIn string) (string, string, error) {
	if err := ValidateInputKey(cliIn); err != nil {
		return "", "", err
	}

	substrings := strings.Split(cliIn, "/")
	if len(substrings) < 2 {
		return "", "", errors.New("Unable to split the input into environment and application name.")
	}
	return substrings[0], substrings[1], nil
}

// ReadSecrets reads key value pairs from file
func ReadSecrets(secretsFile string) (map[string]string, error) {
	file, err := os.Open(secretsFile)
	if err != nil {
		return nil, err
	}
	secrets := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if !strings.HasPrefix(line, "#") {
			keyVal := strings.SplitN(line, "=", 2)
			if len(keyVal) == 2 {
				secrets[keyVal[0]] = keyVal[1]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return secrets, nil
}

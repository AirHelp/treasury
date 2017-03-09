package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateInputKey checks if cli input is valid
func ValidateInputKey(cliIn string) error {
	match, err := regexp.MatchString(`^[a-zA-Z-_]+\/[a-zA-Z-_]+\/[a-zA-Z-_]*$`, cliIn)
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

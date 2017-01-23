package utils

import (
	"errors"
	"regexp"
	"strings"
)

// validateInputKey checks if cli input is valid
func validateInputKey(cliIn string) (bool, error) {
	return regexp.MatchString(`^[a-z]+/[a-z]+/.+$`, cliIn)
}

// FindEnvironmentApplicationName slices cli input into environment and application names
func FindEnvironmentApplicationName(cliIn string) (string, string, error) {
	validIn, err := validateInputKey(cliIn)
	if err != nil {
		return "", "", err
	}
	if !validIn {
		return "", "", errors.New("Invalid input. Please check documentation.")
	}

	substrings := strings.Split(cliIn, "/")
	if len(substrings) < 2 {
		return "", "", errors.New("Unable to split the input into environment and application name.")
	}
	return substrings[0], substrings[1], nil
}

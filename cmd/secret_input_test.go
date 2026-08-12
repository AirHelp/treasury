package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestReadSecretFromStdin(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{name: "plain secret", input: "superSecretPassword", expected: "superSecretPassword"},
		{name: "trailing newline added by echo", input: "superSecretPassword\n", expected: "superSecretPassword"},
		{name: "trailing windows newline", input: "superSecretPassword\r\n", expected: "superSecretPassword"},
		{name: "multiline secret", input: "-----BEGIN KEY-----\nabc\n-----END KEY-----\n", expected: "-----BEGIN KEY-----\nabc\n-----END KEY-----"},
		{name: "significant whitespace is kept", input: "secret with space \n", expected: "secret with space "},
		{name: "empty input", input: "", err: errEmptySecret},
		{name: "only newline", input: "\n", err: errEmptySecret},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.SetIn(strings.NewReader(testCase.input))

			secret, err := readSecret(cmd)
			if !errors.Is(err, testCase.err) {
				t.Fatalf("expected error %v, got %v", testCase.err, err)
			}
			if secret != testCase.expected {
				t.Errorf("expected secret %q, got %q", testCase.expected, secret)
			}
		})
	}
}

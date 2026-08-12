package client_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/AirHelp/treasury/backend"
	"github.com/AirHelp/treasury/client"
	test "github.com/AirHelp/treasury/test/backend"
	"github.com/AirHelp/treasury/types"
)

// webappEntries is what {{ export "test/webapp/" }} resolves to, the secrets of
// the path in the order of their keys
var webappEntries = []string{
	test.ShortKey1 + "=" + test.KeyValueMap[test.Key1],
	test.ShortKey4 + "=" + test.KeyValueMap[test.Key4],
	test.ShortKey2 + "=" + test.KeyValueMap[test.Key2],
}

func TestEnvFile(t *testing.T) {
	tests := []struct {
		name    string
		envFile string
		want    []string
		wantErr bool
	}{
		{
			name: "plain values",
			envFile: `RAILS_ENV=development
API_TOKEN=test
EMPTY=
WITH_EQUAL_SIGNS=a=b=c`,
			want: []string{
				"RAILS_ENV=development",
				"API_TOKEN=test",
				"EMPTY=",
				"WITH_EQUAL_SIGNS=a=b=c",
			},
		},
		{
			name: "quotes are a part of the value, they are not special",
			envFile: `DOUBLE="quoted value"
SINGLE='quoted value'`,
			want: []string{
				`DOUBLE="quoted value"`,
				`SINGLE='quoted value'`,
			},
		},
		{
			name: "comments",
			envFile: `# a comment on its own line

API_TOKEN=test # only for now
	# an indented comment
HASH_INSIDE=pass#word
HASH_FIRST=#word`,
			want: []string{
				"API_TOKEN=test",
				"HASH_INSIDE=pass#word",
				"HASH_FIRST=",
			},
		},
		{
			name:    "single secret",
			envFile: `COCKPIT_API_PASSWORD={{ read "` + test.Key1 + `" }}`,
			want:    []string{"COCKPIT_API_PASSWORD=" + test.KeyValueMap[test.Key1]},
		},
		{
			name:    "single secret with a comment",
			envFile: `COCKPIT_API_PASSWORD={{ read "` + test.Key1 + `" }} # the one we use locally`,
			want:    []string{"COCKPIT_API_PASSWORD=" + test.KeyValueMap[test.Key1]},
		},
		{
			name:    "whole path, variables named after the last part of the key",
			envFile: `{{ export "test/webapp/" }}`,
			want:    webappEntries,
		},
		{
			name:    "whole path with a comment",
			envFile: `{{export "test/webapp/"}} # everything the app needs`,
			want:    webappEntries,
		},
		{
			// the entries are handed to exec.Cmd, which resolves a duplicate key
			// to its last value, so an override is a repetition
			name: "an entry of the same name is repeated, the later one overrides",
			envFile: `{{ export "test/webapp/" }}
` + test.ShortKey1 + `=overridden
API_TOKEN=first
API_TOKEN=second`,
			want: slices.Concat(webappEntries, []string{
				test.ShortKey1 + "=overridden",
				"API_TOKEN=first",
				"API_TOKEN=second",
			}),
		},
		{
			name:    "unknown secret",
			envFile: `PASSWORD={{ read "test/webapp/no_such_key" }}`,
			wantErr: true,
		},
		{
			name:    "invalid secret path",
			envFile: `PASSWORD={{ read "no_such_path" }}`,
			wantErr: true,
		},
		{
			name:    "invalid variable name",
			envFile: `2FAST2FURIOUS=nope`,
			wantErr: true,
		},
		{
			name:    "line which is neither an entry nor a directive",
			envFile: `just a line`,
			wantErr: true,
		},
	}

	treasury := newTreasury(t, &test.MockBackendClient{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := treasury.EnvFile(envFilePath(t, tt.envFile))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Client.EnvFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.EnvFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvFileMissingFile(t *testing.T) {
	treasury := newTreasury(t, &test.MockBackendClient{})

	// the run command relies on this error to tell the user about --env-file
	if _, err := treasury.EnvFile(filepath.Join(t.TempDir(), ".env.treasury")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Client.EnvFile() error = %v, want a %v one", err, fs.ErrNotExist)
	}
}

// TestEnvFileFetches guards what makes the run command usable: a path is read
// once no matter how many entries refer to it, and every remaining key is read
// by name on one batched call, so nothing outside the listed keys is decrypted.
func TestEnvFileFetches(t *testing.T) {
	tests := []struct {
		name    string
		envFile string
		want    []types.GetObjectsInput
	}{
		{
			name: "keys of any number of paths are read by name on one call",
			envFile: `FIRST={{ read "` + test.Key1 + `" }}
SECOND={{ read "` + test.Key2 + `" }}
THIRD={{ read "` + test.Key3 + `" }}
RAILS_ENV=development`,
			want: []types.GetObjectsInput{{Keys: []string{test.Key1, test.Key2, test.Key3}}},
		},
		{
			name: "the same key asked for twice is read once",
			envFile: `FIRST={{ read "` + test.Key1 + `" }}
AGAIN={{ read "` + test.Key1 + `" }}`,
			want: []types.GetObjectsInput{{Keys: []string{test.Key1}}},
		},
		{
			name:    "a whole path takes a call of its own",
			envFile: `{{ export "test/webapp/" }}`,
			want:    []types.GetObjectsInput{{Prefix: "test/webapp/"}},
		},
		{
			name: "a key an exported path already covers is free",
			envFile: `{{ export "test/webapp/" }}
RENAMED={{ read "` + test.Key1 + `" }}`,
			want: []types.GetObjectsInput{{Prefix: "test/webapp/"}},
		},
		{
			name: "paths are read once each, the keys they cover cost nothing",
			envFile: `{{ export "test/webapp/" }}
FIRST={{ read "` + test.Key1 + `" }}
SECOND={{ read "` + test.Key2 + `" }}
THIRD={{ read "` + test.Key3 + `" }}
{{ export "test/cockpit/" }}
{{ export "test/webapp/" }}`,
			want: []types.GetObjectsInput{{Prefix: "test/webapp/"}, {Prefix: "test/cockpit/"}},
		},
		{
			name:    "plain values need no call at all",
			envFile: `RAILS_ENV=development`,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := &countingBackend{}
			if _, err := newTreasury(t, backend).EnvFile(envFilePath(t, tt.envFile)); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(backend.calls, tt.want) {
				t.Errorf("Client.EnvFile() fetched %+v, want %+v", backend.calls, tt.want)
			}
		})
	}
}

// TestEnvFileKeepsSecretsVerbatim covers what the line by line rendering is for:
// a certificate keeps its newlines, its quotes and the hash inside it, none of
// which the plain value rules are allowed to touch
func TestEnvFileKeepsSecretsVerbatim(t *testing.T) {
	const key = "test/webapp/certificate"
	awkward := `"-----BEGIN KEY-----` + "\nline # two\n" + `-----END KEY-----"`

	treasury := newTreasury(t, &singleSecretBackend{key: key, value: awkward})

	got, err := treasury.EnvFile(envFilePath(t, `CERTIFICATE={{ read "`+key+`" }}`))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"CERTIFICATE=" + awkward}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Client.EnvFile() = %q, want %q", got, want)
	}
}

func newTreasury(t *testing.T, backend backend.API) *client.Client {
	t.Helper()
	treasury, err := client.New(&client.Options{Backend: backend})
	if err != nil {
		t.Fatal(err)
	}
	return treasury
}

func envFilePath(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env.treasury")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

// singleSecretBackend serves one secret, whatever the path asked for
type singleSecretBackend struct {
	backend.API
	key, value string
}

func (s *singleSecretBackend) GetObjects(*types.GetObjectsInput) (*types.GetObjectsOuput, error) {
	return &types.GetObjectsOuput{Secrets: map[string]string{s.key: s.value}}, nil
}

// countingBackend records every fetch of secrets from the backend
type countingBackend struct {
	backend.API
	calls []types.GetObjectsInput
}

func (c *countingBackend) GetObjects(input *types.GetObjectsInput) (*types.GetObjectsOuput, error) {
	c.calls = append(c.calls, *input)
	return (&test.MockBackendClient{}).GetObjects(input)
}

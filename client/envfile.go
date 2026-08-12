package client

import (
	"fmt"
	"maps"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"text/template"
)

var variableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// EnvFile resolves an environment file into a list of KEY=VALUE entries ready to
// be handed over to a subprocess. Secrets are kept in memory only.
//
// Every line is a Go template, so the directives are plain template functions,
// the same ones the template command uses:
//
//	{{ export "development/webapp/" }}               all secrets from the path
//	API_PASSWORD={{ read "development/auth/PASS" }}  a single secret
//	API_TOKEN=test                                   a plain value
//
// Entries come back in the order they appear in the file, duplicates included.
// A later entry of the same name overrides an earlier one, which is what
// exec.Cmd does with a duplicate key in Env.
//
// The file is rendered twice. The first pass fetches nothing, it only notes what
// the file asks for, so the secrets can be read in as few calls as possible, and
// so a malformed file fails before a single secret leaves the store. The second
// pass renders the values fetched in between.
func (c *Client) EnvFile(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")

	// the first pass renders no values, so the entries it collects are
	// meaningless and go away with the envFile they were collected into
	var keys, prefixes []string
	if err := (&envFile{}).render(lines, note(&keys), note(&prefixes)); err != nil {
		return nil, err
	}

	env, err := c.fetch(keys, prefixes)
	if err != nil {
		return nil, err
	}
	if err := env.render(lines, env.read, env.export); err != nil {
		return nil, err
	}
	return env.entries, nil
}

// note is a first pass directive: it records what the file asks for, without
// duplicates, and renders nothing
func note(names *[]string) func(string) (string, error) {
	return func(name string) (string, error) {
		if !slices.Contains(*names, name) {
			*names = append(*names, name)
		}
		return "", nil
	}
}

// envFile holds every secret an environment file asked for, whether it came in
// by name or as a part of a whole path, and the entries the file resolves to
type envFile struct {
	client  *Client
	secrets map[string]string
	// paths already read, so none of them is read twice
	paths map[string]bool
	// KEY=VALUE entries, in the order the file lists them
	entries []string
}

// fetch reads everything the file asked for: a call per exported path, and the
// keys left over batched together, so reading a handful of secrets never costs
// more than a handful of names on one call
func (c *Client) fetch(keys, prefixes []string) (*envFile, error) {
	env := &envFile{
		client:  c,
		secrets: make(map[string]string),
		paths:   make(map[string]bool),
	}
	for _, prefix := range prefixes {
		if err := env.walk(prefix); err != nil {
			return nil, err
		}
	}

	// a key an exported path already covers costs nothing extra
	var missing []string
	for _, key := range keys {
		if _, found := env.secrets[key]; !found {
			missing = append(missing, key)
		}
	}
	values, err := c.ReadKeys(missing)
	if err != nil {
		return nil, err
	}
	maps.Copy(env.secrets, values)
	return env, nil
}

// render resolves every line of an environment file with the given directives,
// collecting the entries the file declares
func (e *envFile) render(lines []string, read, export func(string) (string, error)) error {
	directives := template.FuncMap{"read": read, "export": export}

	for index, raw := range lines {
		lineNo := index + 1
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name, value, isEntry := strings.Cut(line, "=")
		if !isEntry {
			// an export directive is the only thing allowed to stand on its own,
			// and it renders nothing
			rendered, err := renderLine(directives, trimComment(line))
			if err != nil {
				return fmt.Errorf("line %d: %w", lineNo, err)
			}
			if rendered != "" {
				return fmt.Errorf("line %d: expected KEY=VALUE or an export directive, got %q", lineNo, line)
			}
			continue
		}

		name = strings.TrimSpace(name)
		if !variableName.MatchString(name) {
			return fmt.Errorf("line %d: %q is not a valid environment variable name", lineNo, name)
		}
		rendered, err := renderLine(directives, trimComment(strings.TrimSpace(value)))
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNo, err)
		}
		e.add(name, rendered)
	}
	return nil
}

// renderLine resolves the directives of a single line. Rendering line by line
// keeps a secret value opaque, whatever it contains, newlines included.
func renderLine(directives template.FuncMap, line string) (string, error) {
	tmpl, err := template.New("env").Funcs(directives).Parse(line)
	if err != nil {
		return "", err
	}
	var rendered strings.Builder
	if err := tmpl.Execute(&rendered, nil); err != nil {
		return "", err
	}
	return rendered.String(), nil
}

func (e *envFile) add(name, value string) {
	e.entries = append(e.entries, name+"="+value)
}

// read resolves a single secret. Everything is fetched by now, unless the file
// took a different branch once the values were known, in which case the secret
// is fetched here.
func (e *envFile) read(key string) (string, error) {
	if value, found := e.secrets[key]; found {
		return value, nil
	}
	values, err := e.client.ReadKeys([]string{key})
	if err != nil {
		return "", err
	}
	value, found := values[key]
	if !found {
		return "", fmt.Errorf("secret %q not found", key)
	}
	e.secrets[key] = value
	return value, nil
}

// export declares an entry for every secret of a path, each named after the last
// part of its key. It renders nothing itself.
func (e *envFile) export(prefix string) (string, error) {
	if err := e.walk(prefix); err != nil {
		return "", err
	}
	prefix = withSlash(prefix)
	for _, key := range slices.Sorted(maps.Keys(e.secrets)) {
		if strings.HasPrefix(key, prefix) {
			e.add(path.Base(key), e.secrets[key])
		}
	}
	return "", nil
}

// walk reads every secret of a path, at most once per path
func (e *envFile) walk(prefix string) error {
	prefix = withSlash(prefix)
	if e.paths[prefix] {
		return nil
	}
	secrets, err := e.client.ReadGroup(prefix)
	if err != nil {
		return err
	}
	for _, secret := range secrets {
		e.secrets[secret.Key] = secret.Value
	}
	e.paths[prefix] = true
	return nil
}

func withSlash(prefix string) string {
	return strings.TrimSuffix(prefix, "/") + "/"
}

// trimComment removes a trailing comment from a value. A comment starts with a
// '#' which opens the value or follows a space or a tab, so a '#' in the middle
// of a value is kept while a value of its own beginning with one is a comment.
// Quotes are not special, they are a part of the value like any other character.
func trimComment(value string) string {
	for i := 0; i < len(value); i++ {
		if value[i] == '#' && (i == 0 || value[i-1] == ' ' || value[i-1] == '\t') {
			return strings.TrimSpace(value[:i])
		}
	}
	return value
}

APPLICATION_SECRET_KEY={{ readFromEnv .Environment "treasury/key2" }}
NAME={{ .Name }}

# export secrets in flexible way
{{ range $key, $value := exportMap "development/treasury/" }}
{{ $key }}={{ $value }}{{ end }}

# export secrets as key=value
{{ export "development/treasury/" }}

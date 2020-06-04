APPLICATION_SECRET_KEY={{ readFromEnv "development" "treasury/key2" }}

# export secrets in flexible way
{{ range $key, $value := exportMap "development/treasury/" }}
{{ $key }}={{ $value }}{{ end }}

# export secrets as key=value
{{ exportFromEnv "development" "treasury" }}

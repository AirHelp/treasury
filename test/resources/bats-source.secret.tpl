APPLICATION_SECRET_KEY={{ read "development/treasury/key2" }}

# export secrets in flexible way
{{ range $key, $value := exportMap "development/treasury/" }}
{{ $key }}={{ $value }}{{ end }}

# export secrets as key=value
{{ export "development/treasury/" }}

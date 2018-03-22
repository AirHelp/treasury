package utils

import "testing"

func TestValidateInputKey(t *testing.T) {
	validTestStrings := []string{
		"staging/webapp/cockpit_api_pass",
		"STAGING/WeBapp/cockpit_api-pass",
		"Integration/claim-score/elasticsearch_url",
		"staging/webapp/1coc3_k--pit123",
		"staging/wordpress-v2/WP_MIXPANEL_API_KEY",
	}
	for _, testString := range validTestStrings {
		if err := ValidateInputKey(testString); err != nil {
			t.Error(err)
		}
	}

	invalidTestStrings := []string{
		"staging/webapp",
		"45678901jbf",
		"asasa/123!12/312313",
		"1231/1231*2/312313",
	}
	for _, testString := range invalidTestStrings {
		if err := ValidateInputKey(testString); err == nil {
			t.Errorf("expected error for test string: %s", testString)
		}
	}
}

func TestValidateInputKeyPattern(t *testing.T) {
	validTestStrings := []string{
		"staging/webapp/cockpit_api_pass",
		"STAGING/WeBapp/cockpit_api-pass",
		"Integration/claim-score/elasticsearch_url",
		"staging/webapp/1coc3_k--pit123",
		"staging/webapp/",
	}
	for _, testString := range validTestStrings {
		if err := ValidateInputKeyPattern(testString); err != nil {
			t.Error(err)
		}
	}

	invalidTestStrings := []string{
		"45678901jbf",
		"asasa/123/12/312313",
		"1231/123!12/312313",
	}
	for _, testString := range invalidTestStrings {
		if err := ValidateInputKeyPattern(testString); err == nil {
			t.Errorf("expected error for test string: %s", testString)
		}
	}
}

func TestFindEnvironmentApplicationName(t *testing.T) {
	var validTest = []struct {
		input       string
		environment string
		application string
	}{
		{"staging/webapp/cockpit_api_pass", "staging", "webapp"},
		{"staging/claim-score/cockpit_api_pass", "staging", "claim-score"},
		{"staging/claim_score/cockpit-api_pass", "staging", "claim_score"},
	}
	for _, test := range validTest {
		env, app, err := FindEnvironmentApplicationName(test.input)
		if err != nil {
			t.Error(err)
		}
		if env != test.environment {
			t.Errorf("Invalid environment name for: %s", test.input)
		}
		if app != test.application {
			t.Errorf("Invalid application name for: %s", test.input)
		}
	}

	invalidTestStrings := []string{
		"stupid string",
		"%/&/@#$%^&*",
		"asdad/asdad1!/1adads",
	}
	for _, testString := range invalidTestStrings {
		if _, _, err := FindEnvironmentApplicationName(testString); err == nil {
			t.Errorf("expected error for %s", testString)
		}
	}
}

func TestReadSecrets(t *testing.T) {
	expected := map[string]string{
		"KEY1": "value@$!#A&*()+-1",
		"KEY2": "value2",
	}
	secrets, err := ReadSecrets("../test/resources/properties.env.test")
	if err != nil {
		t.Error(err)
	}
	if len(expected) != len(secrets) {
		t.Errorf("Wrong found secrets paits, expected:%d. got:%d", len(expected), len(secrets))
	}
	for expectedKey, expectedValue := range expected {
		foundValue := secrets[expectedKey]
		if foundValue != expectedValue {
			t.Errorf("Wrong value for key:%s, expected:%s, got:%s", expectedKey, expectedValue, foundValue)
		}
	}

}

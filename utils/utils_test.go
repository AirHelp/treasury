package utils

import "testing"

func TestValidateInputKey(t *testing.T) {
	testString := "staging/webapp/cockpit_api_pass"
	validKey, err := validateInputKey(testString)
	if err != nil {
		t.Error(err)
	}
	if !validKey {
		t.Errorf("expected true, actual false for %s", testString)
	}

	testString = "staging/webapp"
	falseKey, err := validateInputKey(testString)
	if err != nil {
		t.Error(err)
	}
	if falseKey {
		t.Errorf("expected false, actual true for %s", testString)
	}
}

func TestFindEnvironmentApplicationName(t *testing.T) {
	testString := "staging/webapp/cockpit_api_pass"
	env, app, err := FindEnvironmentApplicationName(testString)
	if err != nil {
		t.Error(err)
	}
	if env != "staging" {
		t.Errorf("Invalid environment name")
	}
	if app != "webapp" {
		t.Errorf("Invalid application name")
	}

	testString = "stupid string"
	_, _, err2 := FindEnvironmentApplicationName(testString)
	if err2 == nil {
		t.Errorf("expected error, actual %s", err)
	}

	testString = "%/&/@#$%^&*"
	_, _, err3 := FindEnvironmentApplicationName(testString)
	if err3 == nil {
		t.Errorf("expected error, actual %s", err)
	}

}

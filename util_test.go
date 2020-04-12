package main

import (
	"os"
	"testing"
)

func TestGetEnvTrue(t *testing.T) {
	testEnvVarValue := "testing"
	testEnvVarKey := "TEST_ENV"
	os.Setenv(testEnvVarKey, testEnvVarValue)

	returnEnv := getEnv(testEnvVarKey)

	if returnEnv != testEnvVarValue {
		t.Errorf("Return was incorect, got: %s, want: %s", returnEnv, testEnvVarValue)
	}

}

func TestGetEnvFalse(t *testing.T) {
	testReturn := getEnv("testing")

	if testReturn != "" {
		t.Error("Value returned, we expected an error")
	}

}

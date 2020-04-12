package main

import (
	"fmt"
	"os"
)

func getEnv(envKey string) (envValue string) {

	envValue, ok := os.LookupEnv(envKey)
	if ok != true {
		fmt.Errorf("please set %s environment variable", envKey)
		return
	}

	return envValue
}

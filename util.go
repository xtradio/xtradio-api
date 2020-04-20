package main

import (
	"fmt"
	"os"
)

func getEnv(envKey string) (envValue string) {

	envValue, ok := os.LookupEnv(envKey)
	if ok != true {
		fmt.Printf("please set %s environment variable", envKey)
		return
	}

	return envValue
}

func songHistory(list []Song, song Song) []Song {

	var data []Song

	if len(list) == 3 {
		list = list[:len(list)-1]
	}

	data = append(data, song)

	for _, v := range list {
		data = append(data, v)
	}

	return data
}

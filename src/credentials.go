package main

import (
	"encoding/json"

	"github.com/c-fandango/rocketchat-term/creds"
)

func getCredentials(cachePath string) (map[string]string, error) {

	fileBytes, err := creds.ReadCache(cachePath)

	outputCreds := make(map[string]string)

	if err == nil {

		err = json.Unmarshal(fileBytes, &outputCreds)

		return outputCreds, err
	}

	if config.host == "" {
		outputCreds["host"] = creds.GetUserInput("Enter host: ", false)
	} else {
		outputCreds["host"] = config.host
	}

	if config.token == "" {
		outputCreds["username"] = creds.GetUserInput("Enter username: ", false)
		outputCreds["password"] = creds.GetUserInput("Enter password: ", true)
	} else {
		outputCreds["token"] = config.token
	}

	return outputCreds, err
}

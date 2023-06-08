package main

import (
	"encoding/json"
	"github.com/c-fandango/rocketchat-term/creds"
)

func getCredentials(cachePath string) (map[string]string, error) {

	fileBytes, err := creds.ReadCache(cachePath)

	outputCreds := make(map[string]string)

	if err != nil {

		outputCreds["host"] = creds.GetUserInput("Enter host: ", false)
		outputCreds["username"] = creds.GetUserInput("Enter username: ", false)
		outputCreds["password"] = creds.GetUserInput("Enter password: ", true)

		return outputCreds, nil
	}

	err = json.Unmarshal(fileBytes, &outputCreds)

	return outputCreds, err
}

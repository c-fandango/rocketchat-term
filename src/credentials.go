package main

import (
	"encoding/json"
	"github.com/c-fandango/rocketchat-term/creds"
	"github.com/c-fandango/rocketchat-term/utils"
)

func getCredentials(cachePath string) (map[string]string, error) {
	var cache cacheSchema
	outputCreds := map[string]string{}

	fileBytes, err := creds.ReadCache(cachePath)

	if err != nil {
		nonSecrets := map[string]string{
			"host":     "",
			"username": "",
		}

		secrets := map[string]string{
			"password": "",
		}

		err = creds.GetCredentials(nonSecrets, secrets)

		if err != nil {
			return outputCreds, err
		}

		return utils.MergeStringMaps(nonSecrets, secrets), nil
	}

	json.Unmarshal(fileBytes, &cache)
	outputCreds["host"] = cache.Host
	outputCreds["token"] = cache.Token

	return outputCreds, nil
}

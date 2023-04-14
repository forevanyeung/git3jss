package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Config struct {
	Server       string
	Username     string
	Password     string
	ScriptsDir   string
	ExtensionDir string
	// SkipSsl			string
}

func LoadConfiguration(file string) (Config, error) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	var config Config
	err = json.NewDecoder(configFile).Decode(&config)

	return config, err
}

type AuthResponse struct {
	Token   string
	Expires string
}

// TODO: save auth token to file
func auth(config Config) (string, time.Time) {
	url := url.URL{
		Scheme: "https",
		Host:   config.Server,
		Path:   "api/v1/auth/token",
	}
	fmt.Printf("Url: %s\n", url.String())

	req, _ := http.NewRequest("POST", url.String(), nil)
	req.Header.Add("accept", "application/json")
	req.SetBasicAuth(config.Username, config.Password)

	res, _ := http.DefaultClient.Do(req) // TODO: capture error
	defer res.Body.Close()

	var response AuthResponse
	err := json.NewDecoder(res.Body).Decode(&response)

	// fmt.Printf("Token: %s\n", response.Token)
	fmt.Printf("Expiration: %s\n", response.Expires)

	// fix for milliseconds less than 3 significant digits
	length := len(response.Expires)
	if length < 24 {
		addSigSpaces := 24 - length
		fmt.Printf("Fixing expiration time due to missing %d significant digits\n", addSigSpaces)

		addIndex := -1
		for i := 1; i <= addSigSpaces; i++ {
			response.Expires = response.Expires[:length+addIndex] + "0" + response.Expires[length+addIndex:]
			addIndex++
		}
	}

	// example 2023-02-21T02:55:00.860Z
	expires, err := time.Parse("2006-01-02T15:04:05.000Z", response.Expires)
	if err != nil {
		log.Fatal(err)
	}

	token := "Bearer " + response.Token

	return token, expires
}

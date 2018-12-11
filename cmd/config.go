package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/mimir-news/pkg/dbutil"
)

// Service metadata.
const (
	ServiceName    = "stock-search"
	ServiceVersion = "1.0"
)

var (
	unsecuredRoutes        = []string{"/health"}
	defaultSearchLimit     = 10
	defaultSuggestionLimit = 5
)

type config struct {
	db                   dbutil.Config
	port                 string
	tokenSecret          string
	tokenVerificationKey string
	adminID              string
}

func getConfig() config {
	tokenSecret := getSecret(mustGetenv("TOKEN_SECRETS_FILE"))

	return config{
		db:                   dbutil.MustGetConfig("DB"),
		tokenSecret:          tokenSecret.Secret,
		tokenVerificationKey: tokenSecret.Key,
		port:                 mustGetenv("SERVICE_PORT"),
		adminID:              mustGetenv("ADMIN_USER_ID"),
	}
}

type secret struct {
	Secret string `json:"secret"`
	Key    string `json:"key"`
}

func getSecret(filename string) secret {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var s secret
	err = json.Unmarshal(content, &s)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("No value for key: %s\n", key)
	}

	return val
}

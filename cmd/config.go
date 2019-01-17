package main

import (
	"log"
	"os"

	"github.com/mimir-news/pkg/httputil/auth"

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
	db             dbutil.Config
	port           string
	JWTCredentials auth.JWTCredentials
	adminID        string
}

func getConfig() config {
	jwtCredentials := getJWTCredentials(mustGetenv("JWT_CREDENTIALS_FILE"))

	return config{
		db:             dbutil.MustGetConfig("DB"),
		JWTCredentials: jwtCredentials,
		port:           mustGetenv("SERVICE_PORT"),
		adminID:        mustGetenv("ADMIN_USER_ID"),
	}
}

func getJWTCredentials(filename string) auth.JWTCredentials {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	credentials, err := auth.ReadJWTCredentials(f)
	if err != nil {
		log.Fatal(err)
	}

	return credentials
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("No value for key: %s\n", key)
	}

	return val
}

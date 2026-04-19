package providers

import (
	"net/http"
	"os"
)

type Provider string

const (
	GOOGLE Provider = "GOOGLE"
	GITHUB Provider = "GITHUB"
)

type Actions interface {
	AuthRedirect(*http.Request) (string, error)
	CallbackRedirect(*http.Request) (string, error)
}

type OAuthOptions map[Provider]Actions

func NewOAuthOptions() (OAuthOptions, error) {

	options := OAuthOptions{}

	googleEnable := os.Getenv("GOOGLE_ENABLE") == "1"

	if !googleEnable {
		return options, OAuthNotConfiguredError
	}

	if googleEnable {
		googleEnvs := []string{"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "GOOGLE_CALLBACK"}

		ok := validateEnvs(googleEnvs)

		if !ok {
			return options, GoogleEnvMissingError
		}

		options[GOOGLE] = NewGoogle()
	}

	return options, nil
}

func validateEnvs(envs []string) bool {
	for _, v := range envs {
		if os.Getenv(v) == "" {
			return false
		}
	}

	return true
}

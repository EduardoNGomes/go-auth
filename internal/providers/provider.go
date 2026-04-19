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
	githubEnable := os.Getenv("GITHUB_ENABLE") == "1"

	if !googleEnable && !githubEnable {
		return options, OAuthNotConfiguredError
	}

	if googleEnable {
		googleEnvs := [3]string{"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "GOOGLE_CALLBACK"}

		ok := validateEnvs(googleEnvs)

		if !ok {
			return options, GoogleEnvMissingError
		}

		options[GOOGLE] = NewGoogle()
	}

	if githubEnable {
		githubEnvs := [3]string{"GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET", "GITHUB_CALLBACK"}

		ok := validateEnvs(githubEnvs)

		if !ok {
			return options, GithuEnvMissingError
		}

		options[GOOGLE] = NewGoogle()
	}
	return options, nil
}

func validateEnvs(envs [3]string) bool {
	for _, v := range envs {
		if os.Getenv(v) == "" {
			return false
		}
	}

	return true
}

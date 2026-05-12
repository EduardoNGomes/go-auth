package providers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Provider string

const (
	GOOGLE Provider = "GOOGLE"
	GITHUB Provider = "GITHUB"
)

type User struct {
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
	Location  string `json:"location"`
}

type Actions interface {
	AuthRedirect(*http.Request) (string, error)
	CallbackRedirect(*http.Request) (string, error)
	createJWTToken(User) (string, error)
	getUser(*http.Client) (User, error)
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

		options[GITHUB] = NewGithub()
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

func authCommon(c *oauth2.Config) (string, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		wrappedErr := fmt.Errorf("error creating code: %w", err)
		return "", wrappedErr
	}

	code := uuid.String()

	url := c.AuthCodeURL(code)

	return url, nil

}

type callbackCommonParams struct {
	c *oauth2.Config
	g Actions
	r *http.Request
}

func callbackCommon(p *callbackCommonParams) (string, error) {
	conf := p.c
	code := p.r.FormValue("code")

	token, err := conf.Exchange(context.Background(), code)

	if err != nil {
		wrappedErr := fmt.Errorf("error on validate code: %w", err)
		return "", wrappedErr
	}

	client := conf.Client(context.Background(), token)

	user, err := p.g.getUser(client)

	if err != nil {
		return "", err
	}

	tokenJWT, err := p.g.createJWTToken(user)

	if err != nil {
		wrappedErr := fmt.Errorf("error on create JWT TOKEN: %w", err)
		return "", wrappedErr
	}

	return tokenJWT, nil
}

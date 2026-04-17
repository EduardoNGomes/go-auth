package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Google struct{}

type GoogleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

type GoogleActions interface {
	AuthRedirect(*http.Request) (string, error)
	CallbackRedirect(*http.Request) (string, error)
}

func NewGoogle() *Google {
	return &Google{}
}

func (g *Google) getConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_CALLBACK"),
		Scopes: []string{
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}

func (g *Google) AuthRedirect(r *http.Request) (string, error) {
	conf := g.getConfig()
	uuid, err := uuid.NewRandom()

	if err != nil {
		wrappedErr := fmt.Errorf("error creating code: %w", err)
		return "", wrappedErr
	}

	code := uuid.String() + "--" + r.Referer()

	url := conf.AuthCodeURL(code)

	return url, nil
}

func (g *Google) CallbackRedirect(r *http.Request) (string, error) {
	conf := g.getConfig()
	code := r.FormValue("code")
	referer := strings.Split(r.FormValue("state"), "--")[1]

	parsed, err := url.Parse(referer)

	if err != nil {
		wrappedErr := fmt.Errorf("error on parser referer: %w", err)
		return "", wrappedErr
	}

	base := parsed.Scheme + "://" + parsed.Host

	token, err := conf.Exchange(context.Background(), code)

	if err != nil {
		wrappedErr := fmt.Errorf("error on validate code: %w", err)
		return "", wrappedErr
	}

	client := conf.Client(context.Background(), token)

	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")

	if err != nil {
		wrappedErr := fmt.Errorf("error on get profile info: %w", err)
		return "", wrappedErr
	}

	defer resp.Body.Close()

	var user GoogleUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		wrappedErr := fmt.Errorf("error on decode user: %w", err)
		return "", wrappedErr
	}

	callbackPath := base + os.Getenv("CALLBACK_PATH")

	return callbackPath, nil
}

package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	AuthRedirect() (string, error)
	CallbackRedirect(r *http.Request) (string, error)
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

func (g *Google) AuthRedirect() (string, error) {
	conf := g.getConfig()
	code, err := uuid.NewRandom()

	if err != nil {
		wrappedErr := fmt.Errorf("error creating code: %w", err)

		fmt.Println(wrappedErr)

		return "", wrappedErr
	}

	url := conf.AuthCodeURL(code.String())

	return url, nil
}

func (g *Google) CallbackRedirect(r *http.Request) (string, error) {
	conf := g.getConfig()
	code := r.FormValue("code")

	token, err := conf.Exchange(context.Background(), code)

	if err != nil {
		wrappedErr := fmt.Errorf("error creating code: %w", err)

		fmt.Println(wrappedErr)

		return "", wrappedErr
	}

	client := conf.Client(context.Background(), token)

	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")

	if err != nil {
		wrappedErr := fmt.Errorf("error creating code: %w", err)

		fmt.Println(wrappedErr)

		return "", wrappedErr
	}

	defer resp.Body.Close()

	var user GoogleUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		wrappedErr := fmt.Errorf("error creating code: %w", err)

		fmt.Println(wrappedErr)

		return "", wrappedErr
	}
	return "/home", nil
}

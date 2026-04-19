package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Github struct {
	config *oauth2.Config
}

type GithubUser struct {
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
	Location  string `json:"location"`
}

func NewGithub() *Github {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_CALLBACK"),
		Scopes: []string{
			"user",
		},
		Endpoint: github.Endpoint,
	}
	return &Github{
		config: conf,
	}
}

func (g *Github) AuthRedirect(r *http.Request) (string, error) {
	conf := g.config

	url, err := authCommon(conf, r)

	if err != nil {
		return "", err

	}

	return url, nil
}

func (g *Github) CallbackRedirect(r *http.Request) (string, error) {
	funcParams := callbackCommonParams{
		c: g.config,
		r: r,
		g: g,
	}

	url, err := callbackCommon(&funcParams)
	if err != nil {
		return "", err

	}

	return url, nil

}

func (g *Github) createJWTToken(user User) (string, error) {
	key := []byte(os.Getenv("SECRET"))

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"name":     user.Name,
			"email":    user.Email,
			"location": user.Location,
			"picture":  user.AvatarUrl,
			"exp":      time.Now().Add(1 * time.Minute).Unix(),
			"iat":      time.Now().Unix(),
		},
	)

	s, err := t.SignedString(key)

	if err != nil {
		return "", err
	}

	return s, nil
}

func (g *Github) getUser(client *http.Client) (User, error) {
	resp, err := client.Get("https://api.github.com/user")

	var userG GithubUser

	if err != nil {
		wrappedErr := fmt.Errorf("error on get profile info: %w", err)
		return User{}, wrappedErr
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return User{}, fmt.Errorf("github user api returned status: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&userG); err != nil {
		wrappedErr := fmt.Errorf("error on decode user: %w", err)
		return User{}, wrappedErr
	}

	return userG.toUser(), nil
}

func (u GithubUser) toUser() User {
	return User{
		Name:      u.Name,
		AvatarUrl: u.AvatarUrl,
		Email:     u.Email,
		Location:  u.Location,
	}
}

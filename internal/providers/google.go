package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Google struct {
	config *oauth2.Config
}

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

func NewGoogle() *Google {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_CALLBACK"),
		Scopes: []string{
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}
	return &Google{
		config: conf,
	}
}

func (g *Google) AuthRedirect(r *http.Request) (string, error) {
	conf := g.config

	url, err := authCommon(conf, r)

	if err != nil {
		return "", err

	}

	return url, nil
}

func (g *Google) CallbackRedirect(r *http.Request) (string, error) {
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

func (g *Google) createJWTToken(user User) (string, error) {
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

func (g *Google) getUser(client *http.Client) (User, error) {
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")

	if err != nil {
		wrappedErr := fmt.Errorf("error on get profile info: %w", err)
		fmt.Println("wrappedErr", wrappedErr)
		return User{}, wrappedErr
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return User{}, fmt.Errorf("google user api returned status: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var userG GoogleUser

	if err := json.NewDecoder(resp.Body).Decode(&userG); err != nil {
		wrappedErr := fmt.Errorf("error on decode user: %w", err)
		return User{}, wrappedErr
	}

	return userG.toUser(), nil
}

func (u GoogleUser) toUser() User {
	return User{
		Name:      u.Name,
		AvatarUrl: u.Picture,
		Email:     u.Email,
		Location:  u.Locale,
	}
}

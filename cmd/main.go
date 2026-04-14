package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "YOUR_REDIRECT_URL",
		Scopes: []string{
			//"https://www.googleapis.com/auth/userinfor.email",
			//"https://www.googleapis.com/auth/userinfor.profile",
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state")
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Handle the exchange code to initiate a transport.
	tok, err := conf.Exchange(oauth2.NoContext, "authorization-code")
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(oauth2.NoContext, tok)
	client.Get("...")
}

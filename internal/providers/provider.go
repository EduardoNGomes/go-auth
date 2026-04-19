package providers

import "net/http"

type Provider string

const (
	GOOGLE Provider = "GOOGLE"
	GITHUB Provider = "GITHUB"
)

type Actions interface {
	AuthRedirect(*http.Request) (string, error)
	CallbackRedirect(*http.Request) (string, error)
}

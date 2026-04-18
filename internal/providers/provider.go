package providers

import "net/http"

type Actions interface {
	AuthRedirect(*http.Request) (string, error)
	CallbackRedirect(*http.Request) (string, error)
}

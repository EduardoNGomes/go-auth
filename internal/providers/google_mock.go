package providers

import (
	"net/http"
)

type GoogleMock struct{}

func NewGoogleMock() *GoogleMock {
	return &GoogleMock{}
}

func (g *GoogleMock) AuthRedirect(r *http.Request) (string, error) {
	return "hc", nil
}

func (g *GoogleMock) CallbackRedirect(r *http.Request) (string, error) {
	return "hc", nil
}

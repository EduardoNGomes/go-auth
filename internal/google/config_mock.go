package google

import (
	"net/http"
)

type GoogleMock struct{}

func NewGoogleMock() *GoogleMock {
	return &GoogleMock{}
}

func (g *GoogleMock) AuthRedirect() (string, error) {
	return "hc", nil
}

func (g *GoogleMock) CallbackRedirect(r *http.Request) (string, error) {
	return "hc", nil
}

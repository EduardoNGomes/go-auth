package providers

import (
	"net/http"
)

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (g *Mock) AuthRedirect(r *http.Request) (string, error) {
	return "hc", nil
}

func (g *Mock) CallbackRedirect(r *http.Request) (string, error) {
	return "hc", nil
}

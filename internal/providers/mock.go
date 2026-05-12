package providers

import (
	"net/http"
)

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (g *Mock) AuthRedirect(r *http.Request, code string) (string, error) {
	return "hc", nil
}

func (g *Mock) CallbackRedirect(r *http.Request) (string, error) {
	return "hc", nil
}

func (g *Mock) createJWTToken(User) (string, error) {
	return "mock", nil

}

func (g *Mock) getUser(*http.Client) (User, error) {
	var user User
	return user, nil
}

package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitbhub.com/eduardongomes/go-auth/internal/routes"

	g "gitbhub.com/eduardongomes/go-auth/internal/google"
)

func TestRoutes(t *testing.T) {
	c := g.NewGoogleMock()
	s, _ := routes.NewServer(c)
	serverMock, _ := routes.NewRoutes(s)

	t.Run("[GET] HC route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/hc", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("Expected -> %d\n Receive ->%d", http.StatusOK, response.Code)
		}

		var body string

		err := json.NewDecoder(response.Body).Decode(&body)

		if err != nil {
			t.Fatal(err)
		}
		bodyMsg := "Im breathing"

		if body != bodyMsg {
			t.Errorf("Expected -> %s\n Receive ->%s", bodyMsg, body)
		}

	})

	t.Run("[GET] Login route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		if response.Code != http.StatusTemporaryRedirect {
			t.Errorf("Expected -> %d\n Receive ->%d", http.StatusTemporaryRedirect, response.Code)
		}
	})

	t.Run("[POST] Login route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/login", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("Expected -> %d\n Receive ->%d", http.StatusNotFound, response.Code)
		}
	})

	t.Run("[GET] Callback route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/callback", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		if response.Code != http.StatusPermanentRedirect {
			t.Errorf("Expected -> %d\n Receive ->%d", http.StatusPermanentRedirect, response.Code)
		}
	})

	t.Run("[POST] Callback route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/callback", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		if response.Code != http.StatusNotFound {
			t.Errorf("Expected -> %d\n Receive ->%d", http.StatusNotFound, response.Code)
		}
	})
}

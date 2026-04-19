package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitbhub.com/eduardongomes/go-auth/internal/routes"

	"gitbhub.com/eduardongomes/go-auth/internal/providers"
)

func TestRoutes(t *testing.T) {
	c := providers.NewMock()
	p := map[providers.Provider]providers.Actions{providers.GOOGLE: c}
	s, _ := routes.NewServer(p)
	serverMock, _ := routes.NewRoutes(s)

	t.Run("[GET] HC route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/hc", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		verifyStatusCode(t, http.StatusOK, response.Code)

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

	t.Run("[GET] home route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/home", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		verifyStatusCode(t, http.StatusOK, response.Code)
	})

	t.Run("[GET] Invalid login route ", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login/invalid", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		verifyStatusCode(t, http.StatusNotFound, response.Code)
	})

	t.Run("[GET GOOGLE] Login route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login/google", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		verifyStatusCode(t, http.StatusTemporaryRedirect, response.Code)
	})

	t.Run("[POST GOOGLE] Login route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/login/google", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		verifyStatusCode(t, http.StatusNotFound, response.Code)
	})

	t.Run("[GET] Invalid callback route ", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/callback/invalid", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		verifyStatusCode(t, http.StatusNotFound, response.Code)
	})

	t.Run("[GET GOOGLE] Callback route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/callback/google", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		verifyStatusCode(t, http.StatusPermanentRedirect, response.Code)

	})

	t.Run("[POST GOOGLE] Callback route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/callback/google", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		verifyStatusCode(t, http.StatusNotFound, response.Code)
	})
}

func verifyStatusCode(t *testing.T, expect, receive int) {
	if receive != expect {
		t.Errorf("Expected -> %d\n Receive ->%d", expect, receive)
	}
}

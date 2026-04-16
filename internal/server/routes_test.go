package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gitbhub.com/eduardongomes/go-auth/internal/server"
)

func TestRoutes(t *testing.T) {
	serverMock := server.NewServer()

	t.Run("[GET] Login route", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		serverMock.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("Expected -> %d\n Receive ->%d", http.StatusOK, response.Code)
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

		if response.Code != http.StatusOK {
			t.Errorf("Expected -> %d\n Receive ->%d", http.StatusOK, response.Code)
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

package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"gitbhub.com/eduardongomes/go-auth/internal/pages"
	"gitbhub.com/eduardongomes/go-auth/internal/providers"
)

type Server struct {
	http.Handler
	template     *template.Template
	oauthOptions providers.OAuthOptions
}

const ContentTypeJSON = "application/json"

var HtmlTemplatePath = pages.GetHtmlTemplate()

func NewServer(options providers.OAuthOptions) (*Server, error) {
	return &Server{
		oauthOptions: options,
	}, nil
}
func NewRoutes(s *Server) (*Server, error) {

	template, err := template.ParseFiles(HtmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem loading template %s %v", HtmlTemplatePath, err)
	}

	s.template = template

	router := http.NewServeMux()

	router.Handle("/hc", http.HandlerFunc(s.healthchecker))
	router.Handle("/home", http.HandlerFunc(s.getPage))
	fileServer := http.FileServer(http.Dir("internal/pages"))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	router.Handle("/login/{provider}", http.HandlerFunc(s.login))
	router.Handle("/callback/{provider}", http.HandlerFunc(s.callback))

	s.Handler = router

	return s, nil
}

func (s *Server) healthchecker(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Im breathing")

		}
	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Not Found")
	}
}

func (s *Server) getPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	s.template.Execute(w, nil)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("provider")

	validation := providerVerify(w, r.Method, path)

	if !validation {
		return
	}

	pathProvider := providers.Provider(strings.ToUpper(path))

	var redirectURL string

	switch pathProvider {
	case providers.GOOGLE:
		{
			url, err := s.oauthOptions[providers.GOOGLE].AuthRedirect(r)

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Failed connect on auth route", http.StatusInternalServerError)
				return
			}
			redirectURL = url

		}
	case providers.GITHUB:
		{
			url, err := s.oauthOptions[providers.GITHUB].AuthRedirect(r)

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Failed connect on auth route", http.StatusInternalServerError)
				return
			}
			redirectURL = url

		}
	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Not Found")
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (s *Server) callback(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("provider")

	validation := providerVerify(w, r.Method, path)

	if !validation {
		return
	}

	var JWT string

	pathProvider := providers.Provider(strings.ToUpper(path))
	switch pathProvider {
	case providers.GOOGLE:
		{

			jwt, err := s.oauthOptions[providers.GOOGLE].CallbackRedirect(r)

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Failed connect on callback route", http.StatusInternalServerError)
			}

			JWT = jwt

		}
	case providers.GITHUB:
		{

			jwt, err := s.oauthOptions[providers.GITHUB].CallbackRedirect(r)

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Failed connect on callback route", http.StatusInternalServerError)
			}

			JWT = jwt

		}
	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Not Found")

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    JWT,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	redirectURL := os.Getenv("REDIRECT_URL")

	http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
}

func providerVerify(w http.ResponseWriter, requestMethod, path string) bool {
	if requestMethod != http.MethodGet || path == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Not Found")
		return false
	}

	return true
}

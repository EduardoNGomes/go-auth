package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	g "gitbhub.com/eduardongomes/go-auth/internal/google"
	"gitbhub.com/eduardongomes/go-auth/internal/pages"
)

type Server struct {
	http.Handler
	template      *template.Template
	googleActions g.GoogleActions
}

const ContentTypeJSON = "application/json"

var HtmlTemplatePath = pages.GetHtmlTemplate()

func NewServer(actions g.GoogleActions) (*Server, error) {
	return &Server{
		googleActions: actions,
	}, nil
}
func NewRoutes(s *Server) (*Server, error) {

	template, err := template.ParseFiles(HtmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem loading template %s %v", HtmlTemplatePath, err)
	}

	s.template = template

	router := http.NewServeMux()

	router.Handle("/home", http.HandlerFunc(s.getPage))
	fileServer := http.FileServer(http.Dir("internal/pages"))

	router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	router.Handle("/login", http.HandlerFunc(s.login))
	router.Handle("/callback", http.HandlerFunc(s.callback))
	router.Handle("/hc", http.HandlerFunc(s.healthchecker))
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
	switch r.Method {
	case http.MethodGet:
		url, err := s.googleActions.AuthRedirect(r)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed connect on auth route", http.StatusInternalServerError)
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Not Found")
	}
}

func (s *Server) callback(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		url, err := s.googleActions.CallbackRedirect(r)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed connect on callback route", http.StatusInternalServerError)
		}

		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Not Found")
	}
}

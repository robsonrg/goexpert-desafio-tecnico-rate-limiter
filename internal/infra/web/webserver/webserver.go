package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	Middlewares   map[string]func(http.Handler) http.Handler
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		Middlewares:   make(map[string]func(http.Handler) http.Handler),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) AddMiddleware(name string, middleware func(http.Handler) http.Handler) {
	s.Middlewares[name] = middleware
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, middleware := range s.Middlewares {
		s.Router.Use(middleware)
	}
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}

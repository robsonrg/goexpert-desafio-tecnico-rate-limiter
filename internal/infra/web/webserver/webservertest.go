package webserver

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServerTest struct {
	Router      chi.Router
	Handlers    map[string]http.HandlerFunc
	Middlewares map[string]func(http.Handler) http.Handler
}

func NewWebServerTest() *WebServerTest {
	return &WebServerTest{
		Router:      chi.NewRouter(),
		Handlers:    make(map[string]http.HandlerFunc),
		Middlewares: make(map[string]func(http.Handler) http.Handler),
	}
}

func (s *WebServerTest) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServerTest) AddMiddleware(name string, middleware func(http.Handler) http.Handler) {
	s.Middlewares[name] = middleware
}

func (s *WebServerTest) Start() *httptest.Server {
	s.Router.Use(middleware.Logger)
	for _, middleware := range s.Middlewares {
		s.Router.Use(middleware)
	}
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}
	return httptest.NewServer(s.Router)
}

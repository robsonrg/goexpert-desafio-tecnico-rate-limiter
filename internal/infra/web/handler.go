package web

import "net/http"

type WebHandler struct {
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

func (h *WebHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sample endpoint response!"))
}

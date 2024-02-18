package home

import (
	"github.com/marcospereirampj/golang-open-telemetry/home/handlers/get"
	"net/http"
)

type (
	handler struct {
		get.Handler
	}
)

func NewHomeServiceHandler() handler {
	return handler{}
}

func (h handler) Pattern() string {
	return "/home"
}

//nolint:errcheck
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello guys!"))
}

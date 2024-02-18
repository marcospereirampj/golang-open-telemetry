package get

import "net/http"

type Handler struct {
}

func (h Handler) Method() string {
	return http.MethodGet
}

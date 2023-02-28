package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterNotFoundHandler(r *mux.Router) {
	r.NotFoundHandler = http.HandlerFunc(notFoundHandlerFunc)
}

func notFoundHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

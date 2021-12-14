package server

import (
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func ServerStart() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	http.ListenAndServe(":3000", r)
}
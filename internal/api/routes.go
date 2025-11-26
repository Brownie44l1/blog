package api

import (
	"github.com/go-chi/chi/v5"
    "net/http"
)

func SetUpRoutes(userHandler *UserHandler, blogHandler *BlogHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/users/register", userHandler.Register)

	r.Post("/blogs", blogHandler.Publish)
	r.Post("/blogs/{id}", blogHandler.GetBlog)

	return r
}
package main

import (
	"github.com/go-chi/chi/v5"
	"go_project/internal/transport"
	"net/http"
)

func main() {

	r := chi.NewRouter()

	r.HandleFunc("/users", transport.UsersHandler) //Найти бы способ использовать один вызов для "/users" и "/users/{id}"
	r.HandleFunc("/users/{id}", transport.UsersIdHandler)

	http.ListenAndServe(":8080", r)
}

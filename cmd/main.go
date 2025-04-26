package main

import (
	"github.com/go-chi/chi/v5"
	"go_project/internal/transport"
	"net/http"
)

func main() {

	//dependency injection (DI) для связывания слоёв
	// db := connectToDB()
	// repo := NewUserRepository(db)
	// service := NewUserService(repo)
	// handler := NewHandler(service)
	//new_end

	r := chi.NewRouter()

	r.HandleFunc("/users", transport.UsersHandler)
	r.HandleFunc("/users/{id}", transport.UsersIdHandler)

	http.ListenAndServe(":8080", r)
}

package main

import (

	// "fmt"
	"github.com/go-chi/chi/v5"
	"go_project/internal/repository"
	"go_project/internal/transport"
	"net/http"
)

func main() {

	//dependency injection (DI) для связывания слоёв
	db := repository.ConnectToDB() //TODO Обрабатывать ли ошибку при отсутствии БД, если и так падает в панику?
	defer db.Close()

	// repo := NewUserRepository(db)
	// service := NewUserService(repo)
	// handler := NewHandler(service)
	//new_end

	r := chi.NewRouter()

	r.Get("/users/{id}", transport.GetAllUsersDB(db))
	r.Get("/users", transport.GetAllUsersDB(db))
	r.Put("/users/{id}", transport.PutIdUserDB(db))
	r.Delete("/users/{id}", transport.DeleteIdUserDB(db))
	r.Post("/users", transport.PostUserDB(db))
	// transport.MasterHandler(r, db)

	http.ListenAndServe(":8080", r)
}

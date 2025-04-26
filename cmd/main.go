package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go_project/internal/repository"
	"go_project/internal/transport"
	"net/http"
)

func main() {

	//dependency injection (DI) для связывания слоёв
	db := repository.ConnectToDB()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT id, name FROM users WHERE age > $1", 18)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	// repo := NewUserRepository(db)
	// service := NewUserService(repo)
	// handler := NewHandler(service)
	//new_end

	r := chi.NewRouter()

	r.HandleFunc("/users", transport.UsersHandler)
	r.HandleFunc("/users/{id}", transport.UsersIdHandler)

	http.ListenAndServe(":8080", r)
}

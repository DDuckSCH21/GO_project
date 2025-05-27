package main

import (
	"context"
	// "fmt"
	"github.com/go-chi/chi/v5"
	"go_project/internal/repository"
	"go_project/internal/transport"
	"net/http"
)

func main() {

	//dependency injection (DI) для связывания слоёв
	db := repository.ConnectToDB() //TODO Обработать ошибку при отсутствии БД
	defer db.Close()

	// fmt.Println("INSERT")
	rows, err := db.Query(context.Background(), "INSERT INTO public.users (id, name,age,is_student) VALUES (2,'SanPusan',21,true);") //TODO Заготовка для postUser

	if err != nil {
		panic(err)
	}
	// fmt.Printf("RESULT INSERT%s\n", rows)
	defer rows.Close()

	// fmt.Println("SELECT")
	rows_2, err := db.Query(context.Background(), "SELECT id, name FROM public.users") //TODO Заготовка для getAllUsers
	if err != nil {
		panic(err)
	}

	for rows_2.Next() {
		var id int
		var name string
		err = rows_2.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		// fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	defer rows_2.Close()

	// repo := NewUserRepository(db)
	// service := NewUserService(repo)
	// handler := NewHandler(service)
	//new_end

	r := chi.NewRouter()

	r.HandleFunc("/users", transport.UsersHandler)
	r.HandleFunc("/users/{id}", transport.UsersIdHandler)

	http.ListenAndServe(":8080", r)
}

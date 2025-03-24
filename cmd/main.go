package main

import (
	"go_project/internal/transport"
	"net/http"
)

func main() {
	http.HandleFunc("/users", transport.UsersHandler)
	http.HandleFunc("/users/", transport.UsersIdHandler)
	http.ListenAndServe(":8080", nil)

}

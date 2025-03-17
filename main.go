package main

import (
	"fmt"
	"net/http"
	"strings"
)

type users struct {
	id       int
	name     string
	lastName string
	age      int
}

var globalDB map[int]users

func pathHandler(r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	fmt.Println(parts)
}

func getUsers(w *http.ResponseWriter) {
	// todo вернуть всех users из globalDB
}

func usersHandler(w http.ResponseWriter, r *http.Request) {

	pathHandler(r)

	switch r.Method {
	case http.MethodGet:
		getUsers(&w)
	case http.MethodPut:
	case http.MethodPost:
	case http.MethodDelete:
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func main() {
	http.HandleFunc("/users", usersHandler)

	http.ListenAndServe(":8080", nil)

}

///users/:id

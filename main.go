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

func pathHandler(r *http.Request) []string {
	parts := strings.Split(r.URL.Path, "/")
	fmt.Println(parts)
	return parts
}

func getUsers(w *http.ResponseWriter, parts []string) { //Возвращает всех users
	//TODO вернуть всех users из globalDB
	//ИЛИ конкретного user по id
}

func putUser(w *http.ResponseWriter, parts []string) { //Обновляет данные по id
	//TODO найти в globalDB и обновить те данные, которые пришли
}

func postUser(w *http.ResponseWriter, parts []string) { //Создает нового user
	//Добавить запись в globalDB, вернуть новый id
	//Как-то надо считать порядок id
}

func deleteUser(w *http.ResponseWriter, parts []string) { //Удаляет user по id
	//берем part[2] и по нему удаляем из globalDB
}

func usersHandler(w http.ResponseWriter, r *http.Request) {

	parts := pathHandler(r)

	switch r.Method {
	case http.MethodGet: //GET /users OR GET /users/:id
		getUsers(&w, parts)
	case http.MethodPut: //PUT /users/:id
		putUser(&w, parts)
	case http.MethodPost: //POST /users
		postUser(&w, parts)
	case http.MethodDelete: //DELETE /users/:id
		deleteUser(&w, parts)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func main() {
	http.HandleFunc("/users", usersHandler)

	http.ListenAndServe(":8080", nil)

}

///users/:id

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	Id   int
	Data map[string]interface{}
}

var globalDB = make(map[int]User)

func findNextKey() int {
	maxKey := math.MinInt
	for num := range globalDB {
		if maxKey < num {
			maxKey = num
		}
	}
	return maxKey + 1
}

func pathHandler(r *http.Request, w *http.ResponseWriter) int {

	parts := strings.Split(r.URL.Path, "/")

	id, _ := strconv.Atoi(parts[1]) //TODO Обработать ошибку
	return id
}

func getUsers(w *http.ResponseWriter, id int) { //Возвращает конкретного user по id
	//TODO вернуть конкретного user по id

}

func getAllUsers(w *http.ResponseWriter) { //Возвращает всех users
	//TODO вернуть всех users из globalDB
	//ИЛИ конкретного user по id
}

func putIdUser(w *http.ResponseWriter, r *http.Request, id int) { //Обновляет данные по id
	//TODO найти в globalDB и обновить те данные, которые пришли
}

func postUser(w *http.ResponseWriter, r *http.Request) { //Добавить новую запись в globalDB, возвращает новый id
	defer r.Body.Close() //

	var user User
	newKey := 1

	err := json.NewDecoder(r.Body).Decode(&user.Data) //пока так
	if err != nil {
		http.Error(*w, "Error Decode JSON", http.StatusBadRequest)
	}

	if len(globalDB) != 0 {
		newKey = findNextKey()
	}

	globalDB[newKey] = user
	fmt.Fprintf(*w, "Add new User id=[%d]", newKey)

	// fmt.Println("globalDB.data=", user.Data)

}

func deleteIdUser(w *http.ResponseWriter, id int) { //Удаляет user по id
	//берем part[2] и по нему удаляем из globalDB
}

func usersIdHandler(w http.ResponseWriter, r *http.Request) {
	idUser := pathHandler(r, &w)

	switch r.Method {
	case http.MethodGet: //GET /users/:id
		getUsers(&w, idUser)
	case http.MethodPut: //PUT /users/:id
		putIdUser(&w, r, idUser)
	case http.MethodDelete: //DELETE /users/:id
		deleteIdUser(&w, idUser)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func usersHandler(w http.ResponseWriter, r *http.Request) {

	parts := pathHandler(r, &w)
	_ = parts

	switch r.Method {
	case http.MethodGet: //GET /users
		getAllUsers(&w)
	case http.MethodPost: //POST /users
		postUser(&w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", usersIdHandler)

	http.ListenAndServe(":8080", nil)

}

///users/:id

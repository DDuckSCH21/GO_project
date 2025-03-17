package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	Id   int
	Data map[string]any
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
	id, _ := strconv.Atoi(parts[len(parts)-1]) //TODO Обработать ошибку //Колхозненько
	return id
}

func getIdUsers(w *http.ResponseWriter, id int) { //Возвращает конкретного user по id
	//TODO вернуть конкретного user по id

	user, ok := globalDB[id]
	if ok {
		fmt.Fprintf(*w, "User ID = %d: %v\n", id, user.Data)
	} else {
		fmt.Fprintf(*w, "User ID =%d not found\n", id)
	}

}

func getAllUsers(w *http.ResponseWriter) { //Возвращает всех users
	var buf bytes.Buffer
	if len(globalDB) != 0 {
		for ind, val := range globalDB {
			fmt.Fprintf(&buf, "User ID = %d: %v\n", ind, val)
		}
		fmt.Fprintln(*w, &buf)
	} else {
		fmt.Fprintln(*w, "No Data")
	}
}

func putIdUser(w *http.ResponseWriter, r *http.Request, id int) { //Обновляет данные по id
	//TODO найти в globalDB и обновить те данные, которые пришли
	defer r.Body.Close()
	var newUser User

	user, ok := globalDB[id]
	if ok {
		err := json.NewDecoder(r.Body).Decode(&newUser.Data)
		if err != nil {
			http.Error(*w, "Error: Decode JSON", http.StatusBadRequest)
			return
		}
		user.Data = newUser.Data
		globalDB[id] = user //Сразу обновить данные в мапе нельзя
	}
}

func postUser(w *http.ResponseWriter, r *http.Request) { //Добавить новую запись в globalDB, возвращает новый id
	defer r.Body.Close()

	var user User
	newKey := 1

	err := json.NewDecoder(r.Body).Decode(&user.Data) //пока так
	if err != nil {
		if err.Error() == "EOF" {
			http.Error(*w, "Error: Empty Request", http.StatusBadRequest)
		} else {
			http.Error(*w, "Error: Decode JSON", http.StatusBadRequest)
		}
		return
	}

	if len(globalDB) != 0 {
		newKey = findNextKey()
	}
	user.Id = newKey
	globalDB[newKey] = user
	fmt.Fprintf(*w, "Add new User id=[%d]", newKey)
}

func deleteIdUser(w *http.ResponseWriter, id int) { //Удаляет user по id
	_, ok := globalDB[id]
	if ok {
		delete(globalDB, id)
	} else {
		fmt.Fprintf(*w, "User ID =%d not found\n", id)
	}
}

func usersIdHandler(w http.ResponseWriter, r *http.Request) {
	idUser := pathHandler(r, &w)

	switch r.Method {
	case http.MethodGet: //GET /users/:id
		getIdUsers(&w, idUser)
	case http.MethodPut: //PUT /users/:id
		putIdUser(&w, r, idUser)
	case http.MethodDelete: //DELETE /users/:id
		deleteIdUser(&w, idUser)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet: //GET /users
		fmt.Println("GET /users")
		getAllUsers(&w)
	case http.MethodPost: //POST /users
		fmt.Println("POST /users")
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

package transport

import (
	"encoding/json"
	"fmt"
	"go_project/internal/global"
	"go_project/internal/models"
	"net/http"
	"strconv"
	"strings"
)

func sendStatus(status int, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	if status != 200 {
		w.WriteHeader(status)
	}
}

func pathHandler(r *http.Request, w http.ResponseWriter) int {
	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		http.Error(w, "Error: Parsing path. strconv.Atoi", http.StatusBadRequest) //TODO Обработать ошибку //Колхозненько
		return -1
	}
	return id
}

func getIdUsers(w http.ResponseWriter, id int) { //Возвращает конкретного user по id
	user, ok := global.DB.Get(id)
	if ok {
		sendStatus(http.StatusOK, w) // 200 - по дефолту отправляется, не надо еще раз это делать
		fmt.Fprintf(w, "User ID = %d: %v\n", id, user.Data)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func getAllUsers(w http.ResponseWriter) { //Возвращает всех users

	if !global.DB.IsEmpty() {
		buf := global.DB.GetAll()
		sendStatus(http.StatusOK, w)
		fmt.Fprintln(w, &buf)
	} else {
		fmt.Fprintln(w, "No Data")
	}
}

func putIdUser(w http.ResponseWriter, r *http.Request, id int) { //Обновляет данные по id
	defer r.Body.Close()
	var newUser models.User

	_, ok := global.DB.Get(id)
	if ok {
		err := json.NewDecoder(r.Body).Decode(&newUser.Data)
		if err != nil {
			http.Error(w, "Error: Decode JSON", http.StatusBadRequest)
			return
		}
		global.DB.Set(id, newUser) //TODO Проверить. Хрень будто какая-то нерабочая
		sendStatus(http.StatusOK, w)
	} else {
		http.Error(w, "Error: id not found", http.StatusNotFound)
	}
}

func postUser(w http.ResponseWriter, r *http.Request) { //Добавить новую запись в global.DB, возвращает новый id
	defer r.Body.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user.Data) //пока так
	if err != nil {
		if err.Error() == "EOF" {
			http.Error(w, "Error: Empty Request", http.StatusBadRequest)
		} else {
			http.Error(w, "Error: Decode JSON", http.StatusBadRequest)
		}
		return
	}
	user.Id = global.DB.GetNewKey()
	global.DB.Set(user.Id, user)
	sendStatus(http.StatusCreated, w)
	fmt.Fprintf(w, "Add new User id=[%d]\n", user.Id)
}

func deleteIdUser(w http.ResponseWriter, id int) { //Удаляет user по id
	ok := global.DB.Del(id)
	if ok {
		sendStatus(http.StatusOK, w)
	} else {
		fmt.Fprintf(w, "User ID = %d not found\n", id)
	}
}

func UsersIdHandler(w http.ResponseWriter, r *http.Request) {
	idUser := pathHandler(r, w)

	switch r.Method {
	case http.MethodGet: //GET /users/:id
		getIdUsers(w, idUser)
	case http.MethodPut: //PUT /users/:id
		putIdUser(w, r, idUser)
	case http.MethodDelete: //DELETE /users/:id
		deleteIdUser(w, idUser)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: //GET /users
		getAllUsers(w)
	case http.MethodPost: //POST /users
		postUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

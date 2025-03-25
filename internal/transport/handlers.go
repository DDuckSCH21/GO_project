package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_project/internal/global"
	"go_project/internal/models"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func getNewKey() int {

	if len(global.DB) != 0 {

		maxKey := math.MinInt
		for num := range global.DB {
			if maxKey < num {
				maxKey = num
			}
		}
		return maxKey + 1
	}
	return 1
}

func pathHandler(r *http.Request, w *http.ResponseWriter) int {
	parts := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(parts[len(parts)-1]) //TODO Обработать ошибку //Колхозненько
	return id
}

func getIdUsers(w *http.ResponseWriter, id int) { //Возвращает конкретного user по id
	//TODO вернуть конкретного user по id

	user, ok := global.DB[id]
	if ok {
		fmt.Fprintf(*w, "User ID = %d: %v\n", id, user.Data)
	} else {
		fmt.Fprintf(*w, "User ID =%d not found\n", id)
	}

}

func getAllUsers(w *http.ResponseWriter) { //Возвращает всех users
	var buf bytes.Buffer
	if len(global.DB) != 0 {
		for ind, val := range global.DB {
			fmt.Fprintf(&buf, "User ID = %d: %v\n", ind, val)
		}
		fmt.Fprintln(*w, &buf)
	} else {
		fmt.Fprintln(*w, "No Data")
	}
}

func putIdUser(w *http.ResponseWriter, r *http.Request, id int) { //Обновляет данные по id
	//TODO найти в global.DB и обновить те данные, которые пришли
	defer r.Body.Close()
	var newUser models.User

	user, ok := global.DB[id]
	if ok {
		err := json.NewDecoder(r.Body).Decode(&newUser.Data)
		if err != nil {
			http.Error(*w, "Error: Decode JSON", http.StatusBadRequest)
			return
		}
		user.Data = newUser.Data
		global.DB[id] = user //Сразу обновить данные в мапе нельзя
	}
}

func postUser(w *http.ResponseWriter, r *http.Request) { //Добавить новую запись в global.DB, возвращает новый id
	defer r.Body.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user.Data) //пока так
	if err != nil {
		if err.Error() == "EOF" {
			http.Error(*w, "Error: Empty Request", http.StatusBadRequest)
		} else {
			http.Error(*w, "Error: Decode JSON", http.StatusBadRequest)
		}
		return
	}

	user.Id = getNewKey()
	global.DB[user.Id] = user
	fmt.Fprintf(*w, "Add new User id=[%d]", user.Id)
}

func deleteIdUser(w *http.ResponseWriter, id int) { //Удаляет user по id
	_, ok := global.DB[id]
	if ok {
		delete(global.DB, id)
	} else {
		fmt.Fprintf(*w, "User ID =%d not found\n", id)
	}
}

func UsersIdHandler(w http.ResponseWriter, r *http.Request) {
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

func UsersHandler(w http.ResponseWriter, r *http.Request) {

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

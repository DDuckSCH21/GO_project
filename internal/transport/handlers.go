package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go_project/internal/global"
	"go_project/internal/models"
	"net/http"
	"strconv"
	"strings"
)

//FOR DATABASE

func getAllUsersDB(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(r.Context(), "SELECT * FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var usersArr []models.User
		for rows.Next() {
			var uTmp models.User
			err := rows.Scan(&uTmp.Id, &uTmp.Name, &uTmp.Age, &uTmp.Is_student)
			if err != nil {
				// fmt.Printf("SQL ERR=%s", err)
				http.Error(w, "Error: Scan SQL all users", http.StatusBadRequest)
				return
			}
			usersArr = append(usersArr, uTmp)
		}
		sendStatus(http.StatusOK, w)
		json.NewEncoder(w).Encode(usersArr) //Сама отправка. Стоит добавить обработку err
	}
}

func getIdUsersDB(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		row := db.QueryRow(r.Context(), "SELECT * FROM users where id = $1", id)

		fmt.Printf("getIdUsersDB WORK\n")
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// defer row.Close()
		var user models.User

		errScan := row.Scan(&user.Id, &user.Name, &user.Age, &user.Is_student)
		if errScan != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		sendStatus(http.StatusOK, w)
		json.NewEncoder(w).Encode(user)
	}
}

func putIdUserDB(db *pgxpool.Pool) http.HandlerFunc { //Обновляет заданную запись, полностью
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var user models.User
		// row := db.QueryRow(r.Context())
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Error: Decode JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		_, errEx := db.Exec(r.Context(), //TODO прочитать про db.Exec
			"UPDATE users u SET name = $1, age = $2, is_student = $3 WHERE id = $4",
			user.Name, user.Age, user.Is_student, id) //Если значений нет - оставляет старые

		if errEx != nil {
			http.Error(w, "Error: DB execution", http.StatusBadRequest)
			return
		}
		// fmt.Println(row.Scan())
		//no rows in result set - даже если апдейтнулось
		sendStatus(http.StatusOK, w)

	}
}

func deleteIdUserDB(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO
	}
}

func postUserDB(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO
	}
}

//END FOR DATABASE

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
	user, ok := global.DB.Get(id) //TODO: Заменить на БД
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

func MasterHandler(r *chi.Mux, db *pgxpool.Pool) {

	// fmt.Println("INSERT")
	rows, err := db.Query(context.Background(), "INSERT INTO public.users (name,age,is_student) VALUES ('SanPusan',21,true);") //TODO Заготовка для postUser

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

	r.Get("/users", getAllUsersDB(db))
	r.Get("/users/{id}", getIdUsersDB(db))
	r.Put("/users/{id}", putIdUserDB(db))

	// r.HandleFunc("/users", UsersHandler)
	// r.HandleFunc("/users/{id}", UsersIdHandler)
}

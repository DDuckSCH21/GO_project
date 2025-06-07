package transport

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go_project/internal/models"
	"net/http"
)

//FOR DATABASE

func GetAllUsersDB(db *pgxpool.Pool) http.HandlerFunc {
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
				http.Error(w, "Error: Scan SQL all users", http.StatusBadRequest)
				return
			}
			usersArr = append(usersArr, uTmp)
		}
		sendStatus(http.StatusOK, w)
		json.NewEncoder(w).Encode(usersArr) //Сама отправка. Стоит добавить обработку err
	}
}

func GetIdUsersDB(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		row := db.QueryRow(r.Context(), "SELECT * FROM users where id = $1", id)

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

func PutIdUserDB(db *pgxpool.Pool) http.HandlerFunc { //Обновляет заданную запись, полностью
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Error: Decode JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		_, errEx := db.Exec(r.Context(),
			"UPDATE users u SET name = $1, age = $2, is_student = $3 WHERE id = $4",
			user.Name, user.Age, user.Is_student, id) //Если значений нет - оставляет старые

		if errEx != nil {
			http.Error(w, "Error: DB UPDATE", http.StatusBadRequest)
			return
		}
		// fmt.Printf("putIdUserDB row=%s\n", row)
		//UPDATE 1 если апдейтнулось и 0 если нет
		sendStatus(http.StatusOK, w)

	}
}

func DeleteIdUserDB(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		defer r.Body.Close()

		_, errEx := db.Exec(r.Context(),
			"DELETE from users WHERE id = $1", id)

		if errEx != nil {
			http.Error(w, "Error: DB DELETE", http.StatusBadRequest)
			return
		}
		// fmt.Printf("deleteIdUserDB res=%s\n", row)
		//DELETE 1 если удалил и 0, если нечего было удалять
		sendStatus(http.StatusOK, w)

	}
}

func PostUserDB(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User
		// var newId int
		// row := db.QueryRow(r.Context())
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Error: Decode JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// rowCount := db.QueryRow(r.Context(), "SELECT max(id)+1 FROM users") //Костыль, если в таблице нет автоинкремента
		// rowCount.Scan(&newId)
		_, errEx := db.Exec(r.Context(),
			"INSERT INTO users (name, age, is_student) VALUES ($1, $2, $3)",
			user.Name, user.Age, user.Is_student)

		if errEx != nil {
			http.Error(w, "Error: DB INSERT", http.StatusBadRequest)
			return
		}
		// fmt.Printf("postUserDB res=%s\n", row)
		//"INSERT 0 1" если заинсертил
		sendStatus(201, w)
		// fmt.Fprintf(w, "Add new User id=[%d]\n", newId)
	}
}

//END FOR DATABASE

func sendStatus(status int, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	if status != 200 {
		w.WriteHeader(status)
	}
}

func MasterHandler(r *chi.Mux, db *pgxpool.Pool) {

}

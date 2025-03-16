package main

import "net/http"

func usersHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc(":8080/users", usersHandler)

	http.HandleFunc(":8080/users/:id", usersHandler)

	http.ListenAndServe(":8080", nil)
}

///users/:id

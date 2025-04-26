package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	// "go_project/internal/models"
	"os"
)

//https://eax.me/golang-pgx/
//https://github.com/jackc/pgx/wiki/Getting-started-with-pgx

type UserRepository struct {
	db *pgxpool.Pool
}

func ConnectToDB() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), "postgres://myuser:mypassword321@localhost:5432/GO_project")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Connected!\n")
	return pool
}

// func (r *UserRepository) Get(id string) (usr models.User, status bool) { //замена global.DB.Get(id)
// }

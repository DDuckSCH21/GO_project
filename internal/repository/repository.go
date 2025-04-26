package repository

import (
	"database/sql"
	"go_project/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

// func (r *UserRepository) Get(id string) (usr models.User, status bool) { //замена global.DB.Get(id)
// }

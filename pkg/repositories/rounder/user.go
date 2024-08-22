package rounder

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"
)

var (
	// ErrUserNotFound is returned when the user is not found.
	ErrUserNotFound = errors.New("user not found")
)

func (r *repository) CreateUser(user *models.User) error {
	return user.Save(r.db)
}

func (r *repository) UserByUsername(username string) (*models.User, error) {
	sqlStmt := `SELECT id FROM user WHERE username = ?`
	var id int
	err := r.db.Get(&id, sqlStmt, username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: %s", ErrUserNotFound, username)
		default:
			return nil, fmt.Errorf("error getting user by username: %w", err)
		}
	}

	return models.UserById(r.db, id)
}

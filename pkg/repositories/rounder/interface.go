package rounder

import "github.com/Jacobbrewer1/golf-stats-tracker/pkg/models"

type Repository interface {
	// CreateUser creates a new user.
	CreateUser(user *models.User) error

	// UserByUsername returns the user with the given username.
	UserByUsername(username string) (*models.User, error)

	// CreateRound creates a new round.
	CreateRound(round *models.Round) error
}

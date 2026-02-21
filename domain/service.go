package domain

import (
	"context"
	"cookeasy/models"
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type Service interface {
	DBStatus(ctx context.Context) (bool, error)
	CreateUser(ctx context.Context, newUser models.SignUp) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
}

type DBService struct {
	db *sql.DB
}

func NewDBService(db *sql.DB) *DBService {
	return &DBService{
		db: db,
	}
}

func (s *DBService) DBStatus(ctx context.Context) (bool, error) {
	if err := s.db.PingContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}

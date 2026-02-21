package domain

import (
	"context"
	"cookeasy/models"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

func (s *DBService) CreateUser(ctx context.Context, newUser models.SignUp) error {
	sqlStmt := `INSERT INTO "users" ("name", "email", "password", "gender", "dob", "levelOfCooking") 
				VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, sqlStmt, newUser.Name, newUser.Email, newUser.Password, newUser.Gender, newUser.DOB, newUser.LevelOfCooking)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// 23505 = unique violation
			if pqErr.Code == "23505" {
				return ErrEmailAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (s *DBService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sqlStmt := `SELECT "ID", "name", "email", "password", "gender", "dob", "levelOfCooking"
				FROM "users" WHERE email = $1`

	var user models.User
	err := s.db.QueryRowContext(ctx, sqlStmt, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Gender, &user.DOB, &user.LevelOfCooking)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *DBService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	sqlStmt := `SELECT "ID", "name", "email", "password", "gender", "dob", "levelOfCooking"
				FROM "users" WHERE "ID" = $1`

	var user models.User
	err := s.db.QueryRowContext(ctx, sqlStmt, userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Gender, &user.DOB, &user.LevelOfCooking)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

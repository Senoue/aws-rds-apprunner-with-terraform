package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Senoue/aws-rds-apprunner-with-terraform/domain/models"
	"github.com/Senoue/aws-rds-apprunner-with-terraform/domain/repositories"
	"github.com/gin-gonic/gin"
)

type authRepository struct {
	db *sql.DB
}

// NewAuthRepository initializes a new AuthRepository.
func NewAuthRepository(db *sql.DB) repository.AuthRepository {
	return &authRepository{db: db}
}

// Login finds a user by email and password.
func (r *authRepository) Login(c *gin.Context, email, password string) (*model.User, error) {
	var user model.User

	query := "SELECT id, username, email FROM users WHERE email = ? AND password_hash = ?"
	err := r.db.QueryRowContext(c, query, email, password).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return &user, nil
}

// Register inserts a new user into the database.
func (r *authRepository) Register(c *gin.Context, user *model.User) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"

	result, err := r.db.ExecContext(c, query, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("no rows affected or error occurred")
	}

	return nil
}

// UserInfo returns the user information.
func (r *authRepository) UserInfo(c *gin.Context, id int) (*model.User, error) {
	var user model.User
	query := "SELECT id, username, email FROM users WHERE id = ?"
	err := r.db.QueryRowContext(c, query, id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return &user, nil
}

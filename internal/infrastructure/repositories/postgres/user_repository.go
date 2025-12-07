package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/EugeneNail/motivatr-app-auth/internal/domain"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create writes passed user domain model to the users table.
// A successful query execution sets user's id to the last inserted value, otherwise returns an error.
// CreatedAt timestamp is set before execution.
func (repo *UserRepository) Create(user *domain.User) error {
	row := repo.db.QueryRow(
		`INSERT INTO users (name, email, password, created_at) VALUES($1, $2, $3, $4) RETURNING id`,
		user.Name, user.Email, user.Password, time.Now(),
	)

	if err := row.Scan(&user.Id); err != nil {
		return fmt.Errorf("executing an SQL query: %w", err)
	}

	return nil
}

// Find searches for a certain user in the users table by its id field.
// A successful search returns a non-nil user domain model. Otherwise, it returns an error.
// If no user is found, it returns a nil instead of a user model and does not produce an error.
func (repo *UserRepository) Find(id int) (*domain.User, error) {
	row := repo.db.QueryRow(`SELECT id, name, email, password, created_at FROM users WHERE id = $1`, id)
	user := domain.User{}

	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("scanning a record into user %d: %w", id, err)
	}

	return &user, nil
}

// FindByEmail searches for a certain user in the users table by its id field.
// A successful search returns a non-nil user domain model. Otherwise, it returns an error.
// If no user is found, it returns a nil instead of a user model and does not produce an error.
func (repo *UserRepository) FindByEmail(email string) (*domain.User, error) {
	row := repo.db.QueryRow(`SELECT id, name, email, password, created_at FROM users WHERE email = $1`, email)
	user := domain.User{}

	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan row into the user with email %s: %w", email, err)
	}

	return &user, nil
}

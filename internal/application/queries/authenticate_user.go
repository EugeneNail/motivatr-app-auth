package queries

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/EugeneNail/motivatr-app-auth/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
)

type AuthenticateUserQuery struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserResult struct {
	Token string
}

type AuthenticateUserHandler struct {
	repository domain.UserRepository
}

func NewAuthenticateUserHandler(repository domain.UserRepository) *AuthenticateUserHandler {
	return &AuthenticateUserHandler{
		repository: repository,
	}
}

func (handler *AuthenticateUserHandler) Handle(query AuthenticateUserQuery) (*AuthenticateUserResult, error) {
	result := AuthenticateUserResult{}

	user, err := handler.repository.FindByEmail(query.Email)
	if err != nil {
		return nil, fmt.Errorf("fetching a user by email: %w", err)
	}

	if user == nil {
		result.Token = ""
		return &result, nil
	}

	hash := sha256.New()
	if _, err := hash.Write([]byte(query.Password + os.Getenv("PASSWORD_SALT"))); err != nil {
		return nil, fmt.Errorf("password hashing: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), hash.Sum(nil))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		result.Token = ""
		return &result, nil
	}

	if err != nil {
		return nil, fmt.Errorf("comparing passwords: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.Id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	result.Token, err = token.SignedString([]byte(os.Getenv("JWT_SALT")))
	if err != nil {
		return nil, fmt.Errorf("creating a signed token string: %w", err)
	}

	return &result, nil
}

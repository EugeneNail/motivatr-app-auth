package commands

import (
	"crypto/sha256"
	"fmt"
	"github.com/EugeneNail/motivatr-app-auth/internal/domain"
	"github.com/EugeneNail/motivatr-lib-common/pkg/validation"
	"github.com/EugeneNail/motivatr-lib-common/pkg/validation/rules"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
)

type CreateUserCommand struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type CreateUserResult struct {
	DomainValidationErrors map[string]string
	Token                  string
}

type CreateUserHandler struct {
	repository domain.UserRepository
}

func NewCreateUserHandler(repository domain.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repository: repository}
}

func (handler *CreateUserHandler) Handle(command CreateUserCommand) (*CreateUserResult, error) {
	result := CreateUserResult{}

	existingUser, err := handler.repository.FindByEmail(command.Email)
	if err != nil {
		return nil, fmt.Errorf("retrieving a user from the DB: %w", err)
	}

	validator := validation.NewValidator(map[string]any{}, map[string][]rules.RuleFunc{})
	if existingUser != nil {
		validator.AddError("email", "The email has already been taken")
		result.DomainValidationErrors = validator.Errors()
		return &result, nil
	}

	hash := sha256.New()
	if _, err := hash.Write([]byte(command.Password + os.Getenv("PASSWORD_SALT"))); err != nil {
		return nil, fmt.Errorf("hashing a password: %w", err)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(hash.Sum(nil), 14)
	if err != nil {
		return nil, fmt.Errorf("encrypting a password: %w", err)
	}

	newUser := domain.User{
		Name:     command.Name,
		Email:    command.Email,
		Password: string(encryptedPassword),
	}

	if err = handler.repository.Create(&newUser); err != nil {
		return nil, fmt.Errorf("writing a user to the DB: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(newUser.Id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	result.Token, err = token.SignedString([]byte(os.Getenv("JWT_SALT")))
	if err != nil {
		return nil, fmt.Errorf("creating a token signed string: %w", err)
	}

	return &result, err
}

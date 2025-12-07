package application

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/EugeneNail/motivatr-app-auth/internal/domain"
	"github.com/EugeneNail/motivatr-lib-common/pkg/validation"
	"github.com/EugeneNail/motivatr-lib-common/pkg/validation/rules"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CreateUserHandler struct {
	repository domain.UserRepository
}

func NewCreateUserHandler(repository domain.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repository: repository}
}

func (handler *CreateUserHandler) Handle(request *http.Request) (int, any) {
	command := CreateUserCommand{}

	if err := json.NewDecoder(request.Body).Decode(&command); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("parsing the input: %w", err)
	}

	validator := validation.NewValidator(map[string]any{
		"name":                 command.Name,
		"email":                command.Email,
		"password":             command.Password,
		"passwordConfirmation": command.PasswordConfirmation,
	}, map[string][]rules.RuleFunc{
		"name":                 {rules.Required(), rules.Regex(rules.Alpha), rules.Min(3), rules.Max(50)},
		"email":                {rules.Required(), rules.Regex(rules.Email), rules.Min(5), rules.Max(100)},
		"password":             {rules.Required(), rules.Password(), rules.Max(50)},
		"passwordConfirmation": {rules.Required(), rules.Same("password")},
	})

	if err := validator.Validate(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("validating the command: %w", err)
	}

	if validator.Failed() {
		return http.StatusUnprocessableEntity, validator.Errors()
	}

	existingUser, err := handler.repository.FindByEmail(command.Email)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("retrieving a user from the DB: %w", err)
	}

	if existingUser != nil {
		validator.AddError("email", "The email has already been taken")
		return http.StatusUnprocessableEntity, validator.Errors()
	}

	hash := sha256.New()
	if _, err := hash.Write([]byte(command.Password + os.Getenv("PASSWORD_SALT"))); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("hashing a password: %w", err)
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(hash.Sum(nil), 14)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("encrypting a password: %w", err)
	}

	newUser := domain.User{
		Name:     command.Name,
		Email:    command.Email,
		Password: string(encryptedPassword),
	}

	if err = handler.repository.Create(&newUser); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("writing a user to the DB: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(newUser.Id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	signedString, err := token.SignedString([]byte(os.Getenv("JWT_SALT")))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("creating signed string: %w", err)
	}

	return http.StatusCreated, CreateUserResult{Token: signedString}
}

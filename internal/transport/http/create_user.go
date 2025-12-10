package http

import (
	"encoding/json"
	"fmt"
	"github.com/EugeneNail/motivatr-app-auth/internal/application/commands"
	"github.com/EugeneNail/motivatr-lib-common/pkg/validation"
	"github.com/EugeneNail/motivatr-lib-common/pkg/validation/rules"
	"net/http"
)

func (handler *Handler) CreateUser(request *http.Request) (int, any) {
	command := commands.CreateUserCommand{}
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

	result, err := handler.createUserHandler.Handle(command)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling a CreateUser command: %w", err)
	}

	if len(result.DomainValidationErrors) > 0 {
		return http.StatusUnprocessableEntity, result.DomainValidationErrors
	}

	return http.StatusCreated, result.Token
}

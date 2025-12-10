package http

import (
	"encoding/json"
	"fmt"
	"github.com/EugeneNail/motivatr-app-auth/internal/application/queries"
	"github.com/EugeneNail/motivatr-lib-common/pkg/validation"
	"github.com/EugeneNail/motivatr-lib-common/pkg/validation/rules"
	"net/http"
)

func (handler *Handler) AuthenticateUser(request *http.Request) (int, any) {
	query := queries.AuthenticateUserQuery{}
	if err := json.NewDecoder(request.Body).Decode(&query); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("parsing a request body: %w", err)
	}

	validator := validation.NewValidator(map[string]any{
		"email":    query.Email,
		"password": query.Password,
	}, map[string][]rules.RuleFunc{
		"email":    {rules.Required(), rules.Regex(rules.Email), rules.Min(5), rules.Max(100)},
		"password": {rules.Required(), rules.Password(), rules.Max(50)},
	})

	if err := validator.Validate(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("validating input: %w", err)
	}

	if validator.Failed() {
		return http.StatusUnprocessableEntity, validator.Errors()
	}

	result, err := handler.authenticateUserHandler.Handle(query)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling an AuthenticateUser query: %w", err)
	}

	if len(result.Token) == 0 {
		validator.AddError("email", "These credentials do not match our records")
		validator.AddError("password", "These credentials do not match our records")
		return http.StatusUnauthorized, validator.Errors()
	}

	return http.StatusOK, result.Token
}

package http

import (
	"github.com/EugeneNail/motivatr-app-auth/internal/application/commands"
	"github.com/EugeneNail/motivatr-app-auth/internal/application/queries"
)

type Handler struct {
	createUserHandler       *commands.CreateUserHandler
	authenticateUserHandler *queries.AuthenticateUserHandler
}

func NewHandler(createUserHandler *commands.CreateUserHandler, authenticateUserHandler *queries.AuthenticateUserHandler) *Handler {
	return &Handler{
		createUserHandler:       createUserHandler,
		authenticateUserHandler: authenticateUserHandler,
	}
}

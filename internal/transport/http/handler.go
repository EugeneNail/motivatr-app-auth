package http

import (
	"github.com/EugeneNail/motivatr-app-auth/internal/application/commands"
)

type Handler struct {
	createUserHandler *commands.CreateUserHandler
}

func NewHandler(createUserHandler *commands.CreateUserHandler) *Handler {
	return &Handler{
		createUserHandler: createUserHandler,
	}
}

package main

import (
	"fmt"
	"github.com/EugeneNail/motivatr-app-auth/internal/application/commands"
	"github.com/EugeneNail/motivatr-app-auth/internal/application/queries"
	"github.com/EugeneNail/motivatr-app-auth/internal/infrastructure/config"
	"github.com/EugeneNail/motivatr-app-auth/internal/infrastructure/repositories/postgres"
	transport "github.com/EugeneNail/motivatr-app-auth/internal/transport/http"
	"github.com/EugeneNail/motivatr-lib-common/pkg/databases"
	middlewares "github.com/EugeneNail/motivatr-lib-common/pkg/middlewares/http"
	"net/http"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(fmt.Errorf("initializing a configuration object: %w", err))
	}

	db, err := databases.ConnectToPostgres(cfg.Db.Host, cfg.Db.Port, cfg.Db.Name, cfg.Db.User, cfg.Db.Password)
	if err != nil {
		panic(fmt.Errorf("connecting to postgres: %w", err))
	}

	userRepository := postgres.NewUserRepository(db)

	createUserHandler := commands.NewCreateUserHandler(userRepository)
	authenticateUserHandler := queries.NewAuthenticateUserHandler(userRepository)

	httpHandler := transport.NewHandler(
		createUserHandler,
		authenticateUserHandler,
	)

	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/auth/create-user", middlewares.WriteJsonResponse(httpHandler.CreateUser))
	router.HandleFunc("POST /api/v1/auth/authenticate-user", middlewares.WriteJsonResponse(httpHandler.AuthenticateUser))

	err = http.ListenAndServe("0.0.0.0:10000", router)
	if err != nil {
		panic(fmt.Errorf("starting the server: %w", err))
	}
}

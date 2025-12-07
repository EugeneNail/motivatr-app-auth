package main

import (
	"fmt"
	"github.com/EugeneNail/motivatr-app-auth/internal/application"
	"github.com/EugeneNail/motivatr-app-auth/internal/infrastructure/config"
	"github.com/EugeneNail/motivatr-app-auth/internal/infrastructure/repositories/postgres"
	"github.com/EugeneNail/motivatr-lib-common/pkg/databases"
	"github.com/EugeneNail/motivatr-lib-common/pkg/middlewares"
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

	createUserHandler := application.NewCreateUserHandler(userRepository)

	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/auth/create-user", middlewares.WriteJsonResponse(createUserHandler))

	err = http.ListenAndServe("0.0.0.0:10000", router)
	if err != nil {
		panic(fmt.Errorf("starting the server: %w", err))
	}
}

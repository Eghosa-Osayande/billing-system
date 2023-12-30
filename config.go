package main

import (
	"blanq_invoice/handlers"
	"blanq_invoice/repository"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ApiConfig struct {
	Repo repository.RepoInterface
	App  *fiber.App
}

func NewApiConfig(repo repository.RepoInterface) *ApiConfig {
	return &ApiConfig{
		Repo: repo,
		App:  fiber.New(),
	}
}

func (config *ApiConfig) SetupRoutes() {

	app := config.App

	app.Use(util.ErrorMessageMiddleware)

	authHandler := handlers.AuthHandler{
		Repo: config.Repo,
	}

	authHandler.RegisterHandlers(app.Group("/auth"))

	businessHandler := handlers.BusinessHandler{
		Repo: config.Repo,
	}
	businessRoute := app.Group("/business")
	businessRoute.Use(util.AuthenticatedUserMiddleware)
	businessHandler.RegisterHandlers(businessRoute)

}

func (config *ApiConfig) StartServer(address string) {

	app := config.App

	if err := app.Listen(address); err != nil {
		log.Fatal("Server failed to start")
	}

}

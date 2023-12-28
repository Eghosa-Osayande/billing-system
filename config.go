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

func NewApiConfig() *ApiConfig {
	return &ApiConfig{
		Repo: repository.NewRepo(),
		App:  fiber.New(),
	}
}

func (config *ApiConfig) Setup(address string) {

	app := config.App

	app.Use(util.ErrorMessageMiddleware)
	
	authHandler := handlers.AuthHandler{
		Repo: config.Repo,
	}

	authHandler.RegisterHandlers(app.Group("/auth"))

	if err := app.Listen(address); err != nil {
		log.Fatal("Server failed to start")
	}
}

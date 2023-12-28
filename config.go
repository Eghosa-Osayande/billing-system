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
	authRouter := app.Group("/auth")

	authHandler := handlers.AuthHandler{
		Repo: config.Repo,
	}

	authRouter.Get("/signup", authHandler.HandleSignup)
	authRouter.Get("/login", authHandler.HandleSignup)
	authRouter.Get("/verifyEmail", authHandler.HandleSignup)
	authRouter.Get("/resendEmailOtp", authHandler.HandleSignup)

	if err := app.Listen(address); err != nil {
		log.Fatal("Server failed to start")
	}
}

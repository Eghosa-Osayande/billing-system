package main

import (
	"blanq_invoice/internal/auth"
	"blanq_invoice/internal/business"
	"blanq_invoice/middlewares"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type ApiConfig struct {
	App *fiber.App
	DB *pgx.Conn
}

type ApiConfigParams struct {
	DB *pgx.Conn
	App *fiber.App

}

func NewApiConfig(params ApiConfigParams) *ApiConfig {
	return &ApiConfig{
		DB: params.DB,
		App: params.App,
	}
}

func (config *ApiConfig) SetupRoutes() {

	app := config.App
	app.Use(util.ErrorMessageMiddleware)

	authHandler := auth.NewAuthHandler(auth.NewAuthRepo(config.DB))
	authHandler.RegisterHandlers(app.Group("/auth"))

	businessHandler := business.NewBusinessHandler(business.NewBusinessRepo(config.DB))
	businessRoute := app.Group("/business")
	businessRoute.Use(middlewares.AuthenticatedUserMiddleware)
	businessHandler.RegisterHandlers(businessRoute)
}

func (config *ApiConfig) StartServer(address string) {

	app := config.App
	if err := app.Listen(address); err != nil {
		log.Fatal("Server failed to start")
	}

}

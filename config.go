package main

import (
	"blanq_invoice/internal/auth"
	"blanq_invoice/internal/business"
	"blanq_invoice/util"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ApiConfig struct {
	App *fiber.App
	DB *sql.DB
}

type ApiConfigParams struct {
	DB *sql.DB
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

	businessHandler := business.NewBusinessHandler()
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

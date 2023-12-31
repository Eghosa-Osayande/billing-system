package main

import (
	"blanq_invoice/database"
	"blanq_invoice/internal/auth"
	"blanq_invoice/internal/business"
	"blanq_invoice/middlewares"
	"blanq_invoice/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5"
)

type ApiConfig struct {
	App *fiber.App
	DB  *pgx.Conn
	DocsAddress string
}

type ApiConfigParams struct {
	DB  *pgx.Conn
	App *fiber.App
	DocsAddress string
}

func NewApiConfig(params ApiConfigParams) *ApiConfig {
	return &ApiConfig{
		DB:  params.DB,
		App: params.App,
		DocsAddress: params.DocsAddress,
	}
}

func (config *ApiConfig) SetupRoutes() {
	db := database.New(config.DB)

	app := config.App
	app.Use(util.ErrorMessageMiddleware)

	app.Get(config.DocsAddress, swagger.New(swagger.Config{
		TryItOutEnabled: false,
		DeepLinking:     false,
		DocExpansion:    "none",
	}))

	authHandler := auth.NewAuthHandler(auth.NewAuthRepo(db))
	authHandler.RegisterHandlers(app.Group("/auth"))

	businessHandler := business.NewBusinessHandler(business.NewBusinessRepo(db))
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

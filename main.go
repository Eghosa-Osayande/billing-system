package main

import (
	"context"
	"log"
	"os"

	"blanq_invoice/database"
	_ "blanq_invoice/docs"
	"blanq_invoice/internal/handlers"
	"blanq_invoice/internal/repos"
	"blanq_invoice/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// @title Blanq Invoice API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if err != nil {
		panic(err)
	}

	dbAdress := os.Getenv("DBURL")
	docsAdress := os.Getenv("DOCSURL")
	x := context.Background()
	conn, err := pgx.Connect(x, dbAdress)
	if err != nil {
		panic(err)
	}

	defer conn.Close(x)

	if err != nil {
		panic(err)
	}
	err = conn.Ping(x)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	db := database.New(conn)
	config := repos.NewApiRepos(repos.ApiReposParams{
		ClientRepo:   repos.NewClientRepo(db),
		BusinessRepo: repos.NewBusinessRepo(db),
		AuthRepo:     repos.NewAuthRepo(db),
		InvoiceRepo:  repos.NewInvoiceRepo(db),
		UserRepo:     repos.NewUserRepo(db),
	})

	middlewares.NewUserMustHaveBusinessMiddleware(config)

	app.Use(middlewares.ErrorMessageMiddleware)

	app.Get(docsAdress, swagger.New(swagger.Config{
		TryItOutEnabled: false,
		DeepLinking:     false,
		DocExpansion:    "none",
	}))

	handlers.NewAuthHandler(config).RegisterHandlers(app)

	handlers.NewBusinessHandler(config).RegisterHandlers(app)

	handlers.NewClientHandler(config).RegisterHandlers(app)

	handlers.NewInvoiceHandler(config).RegisterHandlers(app)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Server failed to start")
	}

}

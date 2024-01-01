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
	"blanq_invoice/util"

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
	})

	app.Use(util.ErrorMessageMiddleware)

	app.Get(docsAdress, swagger.New(swagger.Config{
		TryItOutEnabled: false,
		DeepLinking:     false,
		DocExpansion:    "none",
	}))

	authHandler := handlers.NewAuthHandler(config)
	authHandler.RegisterHandlers(app.Group("/auth"))

	businessHandler := handlers.NewBusinessHandler(config)
	businessRoute := app.Group("/business")
	businessRoute.Use(middlewares.AuthenticatedUserMiddleware)
	businessHandler.RegisterHandlers(businessRoute)

	clientHandler := handlers.NewClientHandler(config)
	clientRoute := app.Group("/clients")
	clientRoute.Use(middlewares.AuthenticatedUserMiddleware)
	clientHandler.RegisterHandlers(clientRoute)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Server failed to start")
	}

}

package main

import (
	"context"
	"log"
	"os"

	"blanq_invoice/database"
	"blanq_invoice/internal/handlers"
	"blanq_invoice/internal/repos"
	"blanq_invoice/middlewares"

	"github.com/gofiber/fiber/v2"
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
		log.Println("Error Loading .env",err)
	}

	port := os.Getenv("PORT")
	
	dbAdress := os.Getenv("DBURL")
	
	x := context.Background()
	conn, err := pgx.Connect(x, dbAdress)
	if err != nil {
		panic(err)
	}

	defer conn.Close(x)

	err = conn.Ping(x)
	if err != nil {
		log.Println("Error pinging database",err)
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


	handlers.NewAuthHandler(config).RegisterHandlers(app)

	handlers.NewBusinessHandler(config).RegisterHandlers(app)

	handlers.NewClientHandler(config).RegisterHandlers(app)

	handlers.NewInvoiceHandler(config).RegisterHandlers(app)

	handlers.NewDashboardHandler(config).RegisterHandlers(app)

	handlers.NewUserHandler(config).RegisterHandlers(app)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Server failed to start",err)
	}

}

package main

import (
	"context"
	"os"

	_ "blanq_invoice/docs"

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

	server := NewApiConfig(ApiConfigParams{
		DB: conn, App: app, DocsAddress: docsAdress})

	server.SetupRoutes()

	server.StartServer(":" + port)

}

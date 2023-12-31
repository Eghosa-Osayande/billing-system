package main

import (
	"database/sql"
	"os"

	_ "blanq_invoice/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
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
	conn, err := sql.Open("postgres", dbAdress)

	if err != nil {
		panic(err)
	}
	app := fiber.New()
	app.Get(docsAdress, swagger.New(swagger.Config{
		TryItOutEnabled: false,
		DeepLinking:     false,
		DocExpansion:    "none",
	}))

	server := NewApiConfig(ApiConfigParams{DB: conn, App: app})

	server.SetupRoutes()

	server.StartServer(":" + port)

}

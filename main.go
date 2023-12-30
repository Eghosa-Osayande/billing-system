package main

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

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
	conn, err := sql.Open("postgres", dbAdress)

	if err != nil {
		panic(err)
	}

	server := NewApiConfig(ApiConfigParams{DB: conn, App: fiber.New()})

	server.SetupRoutes()

	server.StartServer(":" + port)

}

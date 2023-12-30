package main

import (
	"blanq_invoice/repository"
	"log"
	"os"

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
	repo, err := repository.NewPostgresRepo(dbAdress)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	defer repo.Close()

	server := NewApiConfig(repo)

	server.SetupRoutes()

	server.StartServer(":" + port)

}

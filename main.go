package main

import (
	"blanq_invoice/api"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	server := api.NewApiServer()
	port := os.Getenv("PORT")
	server.Setup(port)

}

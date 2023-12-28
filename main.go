package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	server := NewApiConfig()
	port := os.Getenv("PORT")
	server.Setup(":" + port)

}

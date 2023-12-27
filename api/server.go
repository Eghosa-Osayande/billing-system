package api

import (
	"blanq_invoice/database"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ApiServer struct {
	Repo   database.Repository
	
	Router *chi.Mux
}

func NewApiServer() *ApiServer {

	return &ApiServer{
		Repo: database.NewDevRepo(),
		
		Router: chi.NewRouter(),
	}
}

func (server *ApiServer) Setup(port string) {

	app := server.Router
	app.Get("/signup", server.HandleSignup)
	app.Get("/login", server.HandleSignup)
	app.Get("/verifyEmail", server.HandleSignup)
	app.Get("/resendEmailOtp", server.HandleSignup)
	
	hserver := http.Server{
		Addr:    "localhost:" + port,
		Handler: app,
	}

	if err := hserver.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start")
	}
}

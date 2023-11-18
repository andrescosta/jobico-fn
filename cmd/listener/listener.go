package main

import (
	"log"
	"net/http"
	"os"

	"github.com/andrescosta/workflew/internal/listener"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	controler := listener.ListenerController{
		Host: os.Getenv("queue.host"),
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	logger := httplog.NewLogger("listener-log", httplog.Options{
		JSON: true,
	})

	router.Mount("/events", controler.Routes(logger))
	//log.Println("Server listening at", config.Host)
	log.Fatal(http.ListenAndServe(os.Getenv("host"), router))
}

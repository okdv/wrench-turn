package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/okdv/wrench-turn/utils"
)

func main() {
	var err error
	if os.Getenv("GO_ENV") == "production" {
		err = godotenv.Load(".env")
		log.Print("Initializing WrenchTurn production environment...")
	} else {
		err = godotenv.Load(".env.dev")
		log.Print("Initializing WrenchTurn development environment...")
	}
	// load env

	if err != nil {
		log.Fatalf("Unable to load .env file: %v", err)
	}

	// construct new env variables
	os.Setenv("FRONTEND_URL", utils.UrlBuilder(os.Getenv("FRONTEND_PROTO"), os.Getenv("FRONTEND_DOMAIN"), os.Getenv("FRONTEND_PORT")))
	os.Setenv("API_URL", utils.UrlBuilder(os.Getenv("API_PROTO"), os.Getenv("API_DOMAIN"), os.Getenv("API_PORT")))

	// initiate router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// set cors for router
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsHandler.Handler)

	// establish api routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the WrenchTurn API")
	})

	// serve router
	log.Printf("WrenchTurn server listening on port %v", os.Getenv("API_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("API_PORT"), r))
}

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

	"github.com/okdv/wrench-turn/controllers"
	"github.com/okdv/wrench-turn/db"
	"github.com/okdv/wrench-turn/utils"
)

func main() {
	var err error
	// check for env type, default to dev, load env file based on env type
	if os.Getenv("GO_ENV") == "production" {
		err = godotenv.Load(".env")
		log.Print("Initializing WrenchTurn production environment...")
	} else {
		err = godotenv.Load(".env.dev")
		log.Print("Initializing WrenchTurn development environment...")
	}

	if err != nil {
		log.Fatalf("Unable to load .env file: %v", err)
	}

	// connect to db
	err = db.ConnectDatabase()
	if err != nil {
		log.Fatal("Unable to connect to SQLite Database")
		return
	}

	// construct new env variables
	os.Setenv("FRONTEND_URL", utils.UrlBuilder(os.Getenv("FRONTEND_PROTO"), os.Getenv("FRONTEND_DOMAIN"), os.Getenv("FRONTEND_PORT")))
	os.Setenv("API_URL", utils.UrlBuilder(os.Getenv("API_PROTO"), os.Getenv("API_DOMAIN"), os.Getenv("API_PORT")))

	// initiate controllers
	authController := controllers.NewAuthController()

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

	r.Post("/auth", authController.Auth)
	r.Get("/logout", authController.Logout)

	// serve router
	log.Printf("WrenchTurn server listening on port %v", os.Getenv("API_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("API_PORT"), r))
}

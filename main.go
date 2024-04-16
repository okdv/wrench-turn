package main

import (
	"encoding/json"
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
)

func main() {
	var err error
	// check for env type, default to dev, load env file based on env type
	if os.Getenv("GO_ENV") == "production" {
		err = godotenv.Load(".env.production")
		log.Print("Initializing WrenchTurn production environment...")
	} else {
		err = godotenv.Load(".env.development")
		log.Print("Initializing WrenchTurn development environment...")
	}

	if err != nil {
		log.Fatalf("Unable to load .env file: %v", err)
	}

	// connect to db
	_, err = db.ConnectDatabase()
	if err != nil {
		log.Fatal("Unable to connect to SQLite Database")
		return
	}

	// initiate controllers
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	jobController := controllers.NewJobController()
	taskController := controllers.NewTaskController()
	vehicleController := controllers.NewVehicleController()
	alertController := controllers.NewAlertController()
	labelController := controllers.NewLabelController()

	// initiate router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// set cors for router
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("PUBLIC_FRONTEND_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsHandler.Handler)

	// establish api routes
	// general routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the WrenchTurn API")
		return
	})
	r.Get("/env", func(w http.ResponseWriter, r *http.Request) {
		// create slice of PUBLIC env vars
		envVars := map[string]string{
			"PUBLIC_FRONTEND_URL": os.Getenv("PUBLIC_FRONTEND_URL"),
			"PUBLIC_API_URL":      os.Getenv("PUBLIC_API_URL"),
		}
		// convert into json response
		jsonData, err := json.Marshal(envVars)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Unable to convert env to JSON response")
			return
		}
		// respond with json
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})
	// auth routes
	r.Post("/auth", authController.Auth)
	r.Get("/logout", authController.Logout)
	r.Get("/verify", authController.Verify(authController.TestVerify))
	r.Get("/refresh", authController.Verify(authController.Refresh))
	// user routes
	r.Get("/users", userController.ListUsers)
	r.Get("/users/{username}", userController.GetUserByUsername)
	r.Delete("/users/{username}", authController.Verify(userController.DeleteUser))
	r.Post("/users/create", userController.CreateUser)
	r.Post("/users/edit", authController.Verify(userController.EditUser))
	r.Post("/users/updatePassword", authController.Verify(userController.UpdatePassword))
	// job routes
	r.Get("/jobs", jobController.ListJobs)
	r.Get("/jobs/{id:[0-9]+}", jobController.GetJob)
	r.Post("/jobs/{jobId:[0-9]+}/assignLabel/{labelId:[0-9]+}", authController.Verify(jobController.AssignJobLabel))
	r.Post("/jobs/create", authController.Verify(jobController.CreateJob))
	r.Post("/jobs/edit", authController.Verify(jobController.EditJob))
	r.Delete("/jobs/{id:[0-9]+}", authController.Verify(jobController.DeleteJob))
	// task routes
	r.Get("/jobs/{jobId:[0-9]+}/tasks", taskController.ListTasks)
	r.Get("/jobs/{jobId:[0-9]+}/tasks/{taskId:[0-9]+}", taskController.GetTask)
	r.Patch("/jobs/{jobId:[0-9]+}/tasks/{taskId:[0-9]+}/complete", authController.Verify(taskController.MarkComplete))
	r.Post("/jobs/{jobId:[0-9]+}/tasks/create", authController.Verify(taskController.CreateTask))
	r.Post("/jobs/{jobId:[0-9]+}/tasks/edit", authController.Verify(taskController.EditTask))
	r.Delete("/jobs/{jobId:[0-9]+}/tasks/{taskId:[0-9]+}", authController.Verify(taskController.DeleteTask))
	// vehicle routes
	r.Get("/vehicles", vehicleController.ListVehicles)
	r.Get("/vehicles/{id:[0-9]+}", vehicleController.GetVehicle)
	r.Post("/vehicles/create", authController.Verify(vehicleController.CreateVehicle))
	r.Post("/vehicles/edit", authController.Verify(vehicleController.EditVehicle))
	r.Delete("/vehicles/{id:[0-9]+}", authController.Verify(vehicleController.DeleteVehicle))
	// alert routes
	r.Get("/alerts", authController.Verify(alertController.ListAlerts))
	r.Get("/alerts/{id:[0-9]+}", authController.Verify(alertController.GetAlert))
	r.Patch("/alerts/{id:[0-9]+}/read", authController.Verify(alertController.MarkRead))
	r.Post("/alerts/create", authController.Verify(alertController.CreateAlert))
	r.Post("/alerts/edit", authController.Verify(alertController.EditAlert))
	r.Delete("/alerts/{id:[0-9]+}", authController.Verify(alertController.DeleteAlert))
	// label routes
	r.Get("/labels", authController.Verify(labelController.ListLabels))
	r.Get("/labels/{id:[0-9]+}", authController.Verify(labelController.GetLabel))
	r.Post("/labels/create", authController.Verify(labelController.CreateLabel))
	r.Post("/labels/edit", authController.Verify(labelController.EditLabel))
	r.Delete("/labels/{id:[0-9]+}", authController.Verify(labelController.DeleteLabel))
	// serve router
	log.Printf("WrenchTurn server listening on port %v", os.Getenv("PUBLIC_API_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PUBLIC_API_PORT"), r))
}

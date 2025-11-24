package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Bayan2019/ai-hackathon-2025-api/configuration"
	"github.com/Bayan2019/ai-hackathon-2025-api/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/Bayan2019/ai-hackathon-2025-api/docs"
)

// @title AI HACKATHON 2025
// @version 1.0
// @description This is a sample server AI-HACKATHON.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host ai-hackathon-2025-api.onrender.com
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("warning: assuming default configuration. .env unreadable: %v\n", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	// fmt.Println(dbURL)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "superaihackathon"
	}

	err = configuration.Connect2DB(dbURL)
	if err != nil {
		log.Println("DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
		fmt.Println(err.Error())
	}

	if configuration.ApiCfg != nil {
		configuration.ApiCfg.JwtSecret = jwtSecret
	} else {
		fmt.Println("No DATABASE_URL")
		configuration.ApiCfg = &configuration.ApiConfiguration{
			JwtSecret: jwtSecret,
		}
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/", handler)

	router.Get("/swagger/*",
		// httpSwagger.Handler(httpSwagger.URL("http://89.207.252.238:8080/swagger/doc.json")))
		httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	v1Router := chi.NewRouter()

	if configuration.ApiCfg.DB != nil {
		authHandlers := controllers.NewAuthHandlers(*configuration.ApiCfg)

		v1Router.Post("/auth", authHandlers.SignIn)
		v1Router.Patch("/auth", authHandlers.SignInCode)

		usersHandlers := controllers.NewUsersHandlers(*configuration.ApiCfg)

		v1Router.Get("/profile", authHandlers.MiddlewareAuth(usersHandlers.GetProfile))
	}

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 10,
	}

	log.Printf("Serving on: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! Welcome to AI HACKATHON 2025!")
}

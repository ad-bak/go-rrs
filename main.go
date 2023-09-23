package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"p1/internal/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT must be set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/health", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Println("Server running on port " + portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

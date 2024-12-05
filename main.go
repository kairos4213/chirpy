package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/kairos4213/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secretKey      string
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	godotenv.Load()
	pf := os.Getenv("PLATFORM")
	if pf == "" {
		log.Fatal("PLATFORM must be set")
	}

	sk := os.Getenv("SECRET_KEY")
	if sk == "" {
		log.Fatal("SECRET_KEY must be set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error establishing connection to database: %v", err)
	}
	dbQueries := database.New(dbConn)

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       pf,
		secretKey:      sk,
	}

	mux := http.NewServeMux()
	fServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", cfg.middleWareMetricsInc(fServerHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)

	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGet)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

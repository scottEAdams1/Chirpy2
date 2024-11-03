package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/scottEAdams1/Chirpy2/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polka_key      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}

	dbQueries := database.New(db)
	platform := os.Getenv("PLATFORM")

	secret := os.Getenv("SECRET")

	polka_key := os.Getenv("POLKA_KEY")

	const port = "8080"
	apiCfg := apiConfig{
		db:        dbQueries,
		platform:  platform,
		secret:    secret,
		polka_key: polka_key,
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	//Admin
	//Metrics
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	//Reset
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	//API
	//Healthz
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	//Users
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	//Chirps
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	//Login
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	//Refresh
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	//Revoke
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	//Polka
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUser)

	fmt.Printf("Serving on port: %s\n", port)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

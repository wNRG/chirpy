package main

import (
	"log"
	"os"
	"net/http"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/wNRG/chirpy/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	const rootPath = "."
	const fileServerPath = "/app/"
	const port = "8080"
	
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("%s", err)
	}
	
	defer db.Close()
	dbQueries := database.New(db)

	apiCfg := apiConfig {
		dbQueries: dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle(fileServerPath, apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics) 
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Println("server running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}


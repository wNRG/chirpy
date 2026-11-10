package main

import (
	"log"
	"net/http"
)

func main() {
	const rootPath = "."
	const fileServerPath = "/app/"
	const port = "8080"
	
	var apiCfg apiConfig
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


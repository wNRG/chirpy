package main

import (
	"log"
	"net/http"
)

func main() {
	const rootPath = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(rootPath)))

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Println("server running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

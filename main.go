package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nlatham1999/aiapp/server"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

func main() {

	fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	server := &server.Server{}

	r := mux.NewRouter()

	// homepage
	r.HandleFunc("/", server.HomeHandler).Methods("GET")

	// add a route
	r.HandleFunc("/prompt", server.PromptHandler).Methods("POST")

	// add a health check route
	r.HandleFunc("/health", server.PromptHandler).Methods("GET")

	// add a static file route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080", // Address and port for the server to listen on
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting server on http://127.0.0.1:8080/")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}

package main

import (
	"log"
	"net/http"
	"url-shortner/handlers"
	"url-shortner/middleware"
	"url-shortner/storage"
)

func main() {
	mux := http.NewServeMux()
	redisStorage := storage.NewRedisStorage("localhost:6379")
	api := handlers.NewAPI(redisStorage)
	// Protect /shorten endpoint with API key auth
	mux.Handle("/shorten", middleware.APIKeyAuth(http.HandlerFunc(api.ShortenHandler)))
	mux.HandleFunc("/", api.RedirectHandler)
	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

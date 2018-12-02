package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// TODO: add analytics
// TODO: lock google analytics down to profsmallpine

func main() {
	// Setup logger.
	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)

	// Load .env file.
	if err := godotenv.Load(); err != nil {
		panic("could not load env!")
	}

	// Minify assets.
	if os.Getenv("ENVIRONMENT") != "production" {
		if success := minifyAssets(); !success {
			panic("could not minify assets!")
		}
	}

	// Build handler.
	h := handler{Logger: logger}

	// Setup routes.
	router := buildRoutes(h)

	// Run server.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), router); err != nil {
		panic("could not start server!")
	}
}

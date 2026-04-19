package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ecoguard/backend/internal"
)

func main() {
	// Initialize database
	if err := internal.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}

	// Create HTTP server with CORS middleware
	mux := http.NewServeMux()

	// Serve frontend files
	fs := http.FileServer(http.Dir("../../frontend"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "../../frontend/index.html")
	})

	// API endpoints with CORS
	mux.HandleFunc("/signup", corsMiddleware(internal.SignupHandler))
	mux.HandleFunc("/login", corsMiddleware(internal.LoginHandler))
	mux.HandleFunc("/update", corsMiddleware(internal.UpdateUserHandler))
	mux.HandleFunc("/weather", corsMiddleware(internal.WeatherHandler))
	mux.HandleFunc("/detect", corsMiddleware(internal.CropDetectorHandler))
	mux.HandleFunc("/calendar", corsMiddleware(internal.CalendarHandler))

	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Graceful shutdown handling
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		fmt.Println("\nShutting down server...")
		if err := server.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Server close error: %v\n", err)
		}
		if err := internal.CloseDB(); err != nil {
			fmt.Fprintf(os.Stderr, "Database close error: %v\n", err)
		}
		os.Exit(0)
	}()

	fmt.Printf("Server running on http://localhost%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
package main

import (
	"fmt"
	"net/http"

	"ecoguard/backend/internal"
)

func main() {
	// Initialize database
	if err := internal.InitDB(); err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	// Serve frontend files
	fs := http.FileServer(http.Dir("../../frontend")) // adjust if your frontend is inside backend
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/index.html")
	})

	// API endpoints
	http.HandleFunc("/signup", internal.SignupHandler)
	http.HandleFunc("/login", internal.LoginHandler)
	http.HandleFunc("/update", internal.UpdateUserHandler)
	http.HandleFunc("/weather", internal.WeatherHandler)
	http.HandleFunc("/detect", internal.CropDetectorHandler)
	http.HandleFunc("/calendar", internal.CalendarHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
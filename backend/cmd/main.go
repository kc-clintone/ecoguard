package main

import (
	"fmt"
	"net/http"

	"ecoguard/backend/internal"
)

func main() {
	// Load users from file
	if err := internal.LoadUsers(); err != nil {
		fmt.Println("Error loading users:", err)
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

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
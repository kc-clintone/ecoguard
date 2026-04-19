package internal

import (
	"encoding/json"
	"net/http"
)

// SignupHandler handles new user signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Validate request body size
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate input
	if len(newUser.Username) < 3 || len(newUser.Password) < 6 {
		http.Error(w, "Username must be 3+ chars, password 6+ chars", http.StatusBadRequest)
		return
	}

	exists, err := UserExists(newUser.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	if err := InsertUser(newUser); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful"})
}

// LoginHandler authenticates a user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := GetUser(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if user.Password != req.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler updates user data
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updated User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := UpdateUser(updated); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})
}
package internal

import (
	"encoding/json"
	"net/http"
)

// SignupHandler handles new user signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	// Check if user already exists
	for _, u := range Users {
		if u.Username == newUser.Username {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
	}

	Users = append(Users, newUser)
	if err := SaveUsers(); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Signup successful"})
}

// LoginHandler authenticates a user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	for _, u := range Users {
		if u.Username == req.Username && u.Password == req.Password {
			json.NewEncoder(w).Encode(u)
			return
		}
	}

	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
}

// UpdateUserHandler updates user data
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var updated User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	for i, u := range Users {
		if u.Username == updated.Username {
			Users[i] = updated
			if err := SaveUsers(); err != nil {
				http.Error(w, "Failed to save user", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}
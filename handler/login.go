package handlers

import (
	allstruct "dod-backend/All_Struct"
	db "dod-backend/Database"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Role    string `json:"role"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user allstruct.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Use the global DB connection
	query := "INSERT INTO usertbl (name, phoneNumber) VALUES (?, ?)"
	_, err := db.DB.Exec(query, user.Name, user.PhoneNumber)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		fmt.Println("Insert error:", err)
		return
	}

	query = "select role from usertbl where phoneNumber=?"
	var role string
	err = db.DB.Get(&role, query, user.PhoneNumber)
	if err != nil {
		http.Error(w, "Failed to get the user", http.StatusInternalServerError)
		fmt.Println("get error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("User registered successfully"))
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "User registered successfully",
		Role:    role,
	})

	fmt.Println("User saved:", user)
}

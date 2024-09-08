package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

// Login and get a token
// (POST /login)
func (server Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	var body gen.PostLoginJSONBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	type respStruct struct {
		Token string
		User  gen.User
	}
	var resp respStruct

	resp.Token, resp.User, err = server.DB.CreateSession(body, r.Context())
	if err != nil {
		if errors.Is(err, dbModels.ErrUserNotFound) || errors.Is(err, dbModels.ErrInvalidPassword) {
			http.Error(w, "Wrong credentials!", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
		return
	}

	// Set the session token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true, // This makes the cookie inaccessible via JavaScript
		Secure:   true, // Set to true if you're using HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Logout and invalidate token
// (POST /logout)
func (server Server) PostLogout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// _ = json.NewEncoder(w).Encode(resp)
}

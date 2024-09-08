package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/bcrypt"
)

// Get all users
// (GET /users)
func (server Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	var resp []gen.User

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Create a new user
// (POST /users)
func (server Server) PostUsers(w http.ResponseWriter, r *http.Request) {
	var userWithPW struct {
		Pwd string `json:"pwd,omitempty"`
		gen.User
	}

	err := json.NewDecoder(r.Body).Decode(&userWithPW)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(log.Writer()).Encode(userWithPW)
	resp, err := server.DB.CreateUser(userWithPW.User, userWithPW.Pwd, r.Context())

	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Field('C') {
			case "23505":
				http.Error(w, "User already exists.", http.StatusConflict)
			default:
				http.Error(w, "Internal server error."+pgErr.Field('C'), http.StatusInternalServerError)
				log.Printf("Error type: %T, Details: %s", err, err.Error())
			}
		} else if errors.Is(err, dbModels.ErrPasswordNotMachRequirements) {
			http.Error(w, " Password dose not match the requirements.", http.StatusUnprocessableEntity)
		} else if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			http.Error(w, "Password to long.", http.StatusUnprocessableEntity)
		} else {
			log.Printf("Error type: %T, Details: %s", err, err.Error())
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Delete a user by ID
// (DELETE /users/{userId})
func (server Server) DeleteUsersUserId(w http.ResponseWriter, r *http.Request, userId string) {
	var resp []gen.User

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Get a user by ID
// (GET /users/{userId})
func (server Server) GetUsersUserId(w http.ResponseWriter, r *http.Request, userId string) {
	var resp []gen.User

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Update a user by ID
// (PUT /users/{userId})
func (server Server) PutUsersUserId(w http.ResponseWriter, r *http.Request, userId string) {
	var resp []gen.User

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Returns currently logged in user.
// (GET /currentUser)
func (server Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token[:7] == "Bearer " {
		token = token[7:]
		user, err := server.DB.VerifySession(token, r.Context())
		if err != nil {
			if errors.Is(err, dbModels.ErrInvalidPassword) || errors.Is(err, dbModels.ErrUserNotFound) {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			var jwterr *jwt.ValidationError
			if errors.As(err, &jwterr) {
				errCode := jwterr.Errors
				// Group 1: Malformed, Unverifiable, Signature Invalid
				if errCode&jwt.ValidationErrorMalformed != 0 ||
					errCode&jwt.ValidationErrorUnverifiable != 0 ||
					errCode&jwt.ValidationErrorSignatureInvalid != 0 {
					http.Error(w, "Token malformed or Signature Invalid.", http.StatusBadRequest)
					return
				} else if errCode&jwt.ValidationErrorExpired != 0 ||
					errCode&jwt.ValidationErrorNotValidYet != 0 {
					http.Error(w, "Token currently not Valid.", http.StatusUnauthorized)
					return
				} else if errCode&jwt.ValidationErrorId != 0 ||
					errCode&jwt.ValidationErrorIssuedAt != 0 ||
					errCode&jwt.ValidationErrorIssuer != 0 ||
					errCode&jwt.ValidationErrorClaimsInvalid != 0 {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				} else {
					http.Error(w, "Malformed Authorization", http.StatusBadRequest)
				}
			}
			http.Error(w, "Internal server error."+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
		return
	}
	http.Error(w, "Malformed Authorization", http.StatusBadRequest)
}

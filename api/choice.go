package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
)

// Get choices by userId
// (GET /users/{userId}/choices)
func (server Server) GetUsersUserIdChoices(w http.ResponseWriter, r *http.Request, userId int) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	if userId == -1 {
		userId = *user.Id
	}
	if user.Role == nil || *user.Role != gen.UserRoleAdmin {
		if userId != *user.Id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
	choices, err := server.DB.GetChoicesByUserId(userId, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(choices)
}

// Get a choice by userId and choiceId
// (GET /users/{userId}/choices/{choiceId})
func (server Server) GetUsersUserIdChoicesChoiceId(w http.ResponseWriter, r *http.Request, userId int, choiceId int) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	if userId == -1 {
		userId = *user.Id
	}
	if user.Role == nil || *user.Role != gen.UserRoleAdmin {
		if userId != *user.Id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
	choice, err := server.DB.GetChoiceByUserIdAndChoiceId(userId, choiceId, r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Choice not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(choice)
}

// Modify or create a choice by userId and choiceId
// (POST /users/{userId}/choices/{choiceId})
func (server Server) PostUsersUserIdChoicesChoiceId(w http.ResponseWriter, r *http.Request, userId int, choiceId int) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	if userId == -1 {
		userId = *user.Id
	}
	if user.Role == nil || *user.Role != gen.UserRoleAdmin {
		if userId != *user.Id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	var choice gen.Choice
	err = json.NewDecoder(r.Body).Decode(&choice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	choice, err = server.DB.CreateOrUpdateChoice(userId, choiceId, choice, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(choice)
	w.WriteHeader(http.StatusCreated)
}

package api

import (
	"net/http"
)

// Get choices by userId
// (GET /users/{userId}/choices)
func (server Server) GetUsersUserIdChoices(w http.ResponseWriter, r *http.Request, userId int) {
	w.WriteHeader(http.StatusNotFound)
}

// Get a choice by userId and choiceId
// (GET /users/{userId}/choices/{choiceId})
func (server Server) GetUsersUserIdChoicesChoiceId(w http.ResponseWriter, r *http.Request, userId string, choiceId int) {

	w.WriteHeader(http.StatusNotFound)
}

// Modify or create a choice by userId and choiceId
// (POST /users/{userId}/choices/{choiceId})
func (server Server) PostUsersUserIdChoicesChoiceId(w http.ResponseWriter, r *http.Request, userId int, choiceId int) {

	w.WriteHeader(http.StatusNotFound)
}

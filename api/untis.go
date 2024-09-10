package api

import (
	"encoding/json"
	"net/http"
)

func (server Server) GetUntisClasses(w http.ResponseWriter, r *http.Request) {
	_, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	classes, err := server.DB.GetClasses(r.Context())
	if err != nil {

		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(classes)
}

func (server Server) GetUntisRooms(w http.ResponseWriter, r *http.Request) {
	_, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	rooms, err := server.DB.GetRooms(r.Context())
	if err != nil {

		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rooms)
}

func (server Server) GetUntisSubjects(w http.ResponseWriter, r *http.Request) {
	_, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	subjects, err := server.DB.GetSubjects(r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(subjects)
}

func (server Server) GetUntisTeachers(w http.ResponseWriter, r *http.Request) {
	_, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	teachers, err := server.DB.GetTeachers(r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(teachers)
}

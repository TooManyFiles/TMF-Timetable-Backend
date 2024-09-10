package api

import (
	"encoding/json"
	"net/http"
)

func (server Server) GetUntisClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := server.DB.GetClasses(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(classes)
}

func (server Server) GetUntisRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := server.DB.GetRooms(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rooms)
}

func (server Server) GetUntisSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := server.DB.GetSubjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(subjects)
}

func (server Server) GetUntisTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := server.DB.GetTeachers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(teachers)
}

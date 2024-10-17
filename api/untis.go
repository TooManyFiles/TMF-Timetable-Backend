package api

import (
	"encoding/json"
	"net/http"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
)

func (server Server) GetUntisClasses(w http.ResponseWriter, r *http.Request) {
	_, _, err := server.isLoggedIn(w, r)
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
	_, _, err := server.isLoggedIn(w, r)
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
	_, _, err := server.isLoggedIn(w, r)
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
	_, _, err := server.isLoggedIn(w, r)
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
func (server Server) GetUntisFetch(w http.ResponseWriter, r *http.Request) {
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		return
	}
	if user.Role == nil || *user.Role != gen.UserRoleAdmin {
		http.Error(w, "Insufficient permission.", http.StatusForbidden)
		return
	}
	err = server.DB.FetchTeachers(r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	err = server.DB.FetchRooms(r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	err = server.DB.FetchSubjects(r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	err = server.DB.FetchClasses(r.Context())

	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

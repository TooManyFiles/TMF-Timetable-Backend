package api

import (
	"encoding/json"
	"net/http"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

// Get events by a user
// (PUT /view)
func (server Server) PutView(w http.ResponseWriter, r *http.Request, params gen.PutViewParams) {
	resp, err := server.DB.GetLesson(dbModels.LessonFilter{}, r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Get events of a week by a user
// (PUT /view/user/{userId})
func (server Server) PutViewUserUserId(w http.ResponseWriter, r *http.Request, userId int, params gen.PutViewUserUserIdParams) {
	resp, err := server.DB.GetLesson(dbModels.LessonFilter{}, r.Context())
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

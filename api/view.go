package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

// Get events by a user
// (PUT /view)
func (server Server) PutView(w http.ResponseWriter, r *http.Request, params gen.PutViewParams) {
	if *params.Duration < 1 || *params.Duration > 7 {
		http.Error(w, "Invalid request body. Duration out of bounce.", http.StatusBadRequest)
		return
	}
	var body gen.PutViewJSONBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}
	user, _, err := server.isLoggedIn(w, r)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	untis_pwd, err := server.DB.GetUntisLoginByHeader(r.Header.Get("Authorization"), r.Context())
	if err != nil {
		http.Error(w, "Internal server error."+err.Error(), http.StatusInternalServerError)
		return
	}
	startdate := time.Now().Truncate(24 * time.Hour)
	if params.Date != nil && !params.Date.IsZero() {
		startdate = params.Date.Time
	}
	enddate := time.Time(startdate)
	enddate = enddate.AddDate(0, 0, *params.Duration)
	for _, classId := range *user.Classes {
		server.DB.FetchLesson(user, untis_pwd, classId, startdate, enddate, r.Context()) //TODO: change 4454
	}

	lessonFilter := dbModels.LessonFilter{
		User:      (&dbModels.User{}).FromGen(user),
		StartDate: startdate,
		EndDate:   enddate,
	}
	if body.Untis.Choice != nil {
		lessonFilter.Choice = (&dbModels.Choice{}).FromGen(*body.Untis.Choice)
	}
	resp, err := server.DB.GetLesson(lessonFilter, r.Context())
	if err != nil {
		http.Error(w, "Internal server error."+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Get events of a week by a user
// (PUT /view/user/{userId})
func (server Server) PutViewUserUserId(w http.ResponseWriter, r *http.Request, userId int, params gen.PutViewUserUserIdParams) {
	if *params.Duration < 1 || *params.Duration > 7 {
		http.Error(w, "Invalid request body. Duration out of bounce.", http.StatusBadRequest)
		return
	}
	var body gen.PutViewUserUserIdJSONBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}
	user, _, err := server.isLoggedIn(w, r)
	if *user.Id != userId && *user.Role != gen.UserRoleAdmin {
		http.Error(w, "Unauthorized to access this data.", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	startdate := time.Now().Truncate(24 * time.Hour)
	if params.Date != nil && !params.Date.IsZero() {
		startdate = params.Date.Time
	}
	enddate := time.Time(startdate)
	enddate = enddate.AddDate(0, 0, *params.Duration)

	lessonFilter := dbModels.LessonFilter{
		User:      dbModels.User{Id: userId},
		StartDate: startdate,
		EndDate:   enddate,
	}
	if body.Untis.Choice != nil {
		lessonFilter.Choice = (&dbModels.Choice{}).FromGen(*body.Untis.Choice)
	}
	resp, err := server.DB.GetLesson(lessonFilter, r.Context())
	if err != nil {
		http.Error(w, "Internal server error."+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

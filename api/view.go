package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
)

type ViewOutput struct {
	Untis     interface{} `json:untis`
	Cafeteria interface{} `json:cafeteria`
	Week      interface{} `json:week`
}

// Get events by a user
// (PUT /view)
func (server Server) PutView(w http.ResponseWriter, r *http.Request, params gen.PutViewParams) {
	if params.Duration == nil {
		params.Duration = new(int)
		*params.Duration = 1
	} else if *params.Duration < 1 || *params.Duration > 7 {
		http.Error(w, "Invalid request body. Duration out of bounce.", http.StatusBadRequest)
		return
	}
	var body gen.PutViewJSONBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}
	user, claims, err := server.isLoggedIn(w, r)
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

	out := ViewOutput{}
	for _, provider := range body.Provider {
		switch provider {
		case gen.PutViewJSONBodyProviderUntis:
			lessons, err := server.UntisView(user, claims, *body.Untis, startdate, enddate, true, r.Context())
			if err != nil {
				http.Error(w, "Internal server error."+err.Error(), http.StatusInternalServerError)
				return
			}
			out.Untis = lessons
		default:
			// Optionally handle unsupported providers
			// http.Error(w, "Unsupported provider: "+string(provider), http.StatusBadRequest)
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(out)
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

	out := ViewOutput{}
	for _, provider := range body.Provider {
		switch provider {
		case gen.PutViewUserUserIdJSONBodyProviderUntis:
			lessons, err := server.UntisView(user, nil, *body.Untis, startdate, enddate, true, r.Context())
			http.Error(w, "Internal server error."+err.Error(), http.StatusInternalServerError)
			out.Untis = lessons
		default:
			// Optionally handle unsupported providers
			http.Error(w, "Unsupported provider: "+string(provider), http.StatusBadRequest)
		}

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

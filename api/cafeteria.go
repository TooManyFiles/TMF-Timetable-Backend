package api

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
)

// GetCafeteria handles the GET request for fetching cafeteria menu
func (server Server) GetCafeteria(w http.ResponseWriter, r *http.Request, params gen.GetCafeteriaParams) {
	var menus []gen.Menu
	var date time.Time
	var days int = 1
	// Handle Date and Duration parameters
	if params.Date != nil {
		// Fetch the menu for the specific date
		date = time.Time(params.Date.Time) // Convert openapi_types.Date to time.Time

	} else {
		date = time.Now()
	}

	if params.Duration != nil {
		days = int(math.Max(1, math.Min(float64(*params.Duration), 7)))
	}
	menus, err := server.DB.FetchMenuForDate(date, days, r.Context())
	if err != nil {
		http.Error(w, "Error fetching menu", http.StatusInternalServerError)
		return
	}
	if len(menus) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Encode the response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(menus)
}

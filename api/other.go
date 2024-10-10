package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Gets the subtile for the week the date is include.
// (GET /week/{date})
func (server Server) GetWeekDate(w http.ResponseWriter, r *http.Request, date openapi_types.Date) {
	if date.Time.IsZero() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode("No date provided.")
	}
	subtitle, err := server.DB.GetWeekSubtitle(date.Time, r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Internal error")
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(subtitle)
}

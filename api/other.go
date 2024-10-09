package api

import (
	"encoding/json"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Gets the subtile for the week the date is include.
// (GET /week/{date})
func (server Server) GetWeekDate(w http.ResponseWriter, r *http.Request, date openapi_types.Date) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode("Not yet implemented")
}

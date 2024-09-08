package api

import (
	"encoding/json"
	"net/http"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
)

// (GET /ping)
func (Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	resp := gen.User{}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

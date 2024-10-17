package api

import (
	"context"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
)

func (server Server) CafeteriaView(date time.Time, days int, ctx context.Context) ([]gen.Menu, error) {
	// Fetch menu data from database
	menus, err := server.DB.FetchMenuForDate(date, days, ctx)
	return menus, err
}

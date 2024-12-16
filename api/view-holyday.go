package api

import (
	"context"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
)

func (server Server) HolydayView(date time.Time, days int, ctx context.Context) ([]gen.Holyday, error) {
	// Fetch menu data from database
	holidays, err := server.DB.GetHolydays(date, days, ctx)
	return holidays, err
}

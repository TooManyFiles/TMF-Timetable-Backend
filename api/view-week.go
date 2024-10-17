package api

import (
	"context"
	"time"
)

func (server Server) WeekView(date time.Time, ctx context.Context) (string, error) {
	// Fetch menu data from database
	subtitle, err := server.DB.GetWeekSubtitle(date, ctx)
	return subtitle, err
}

package googleapi

import (
	"context"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalenderAPI struct {
	ApiKey     string
	CalendarID string
}

// Replace with the specific date in YYYY-MM-DD format.
func (api *GoogleCalenderAPI) GetEvents(startDate time.Time, endDate time.Time) ([]*calendar.Event, error) {

	startDate = startDate.Truncate(24 * time.Hour)
	endDate = endDate.Truncate(24 * time.Hour)
	endDate.Add(time.Hour * 24)

	// Convert the date to RFC3339 format.
	startTime := startDate.Format(time.RFC3339)
	endTime := endDate.Format(time.RFC3339)

	// Create a new Calendar service using the API key.
	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithAPIKey(api.ApiKey))
	if err != nil {
		return nil, err
	}

	// Fetch events for the specified date range.
	events, err := srv.Events.List(api.CalendarID).
		TimeMin(startTime).TimeMax(endTime).
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, err
	}
	return events.Items, nil
}

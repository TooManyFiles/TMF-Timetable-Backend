package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

func getMonday(date time.Time) time.Time {
	// Calculate the offset to the most recent Monday
	offset := int(time.Monday - date.Weekday())
	if offset > 0 {
		offset = -6 // Correct offset when the date is a Sunday
	}

	// Subtract the offset to get the Monday date
	monday := date.AddDate(0, 0, offset)
	return monday.Truncate(24 * time.Hour)
}

func (database *Database) CreateWeekSubtitle(date time.Time, subtitle string, ctx context.Context) error {
	if date.IsZero() {
		return errors.New("date is required")
	}
	monday := getMonday(date)
	dbWeek := dbModels.WeekSubtitle{
		Date:     monday,
		Subtitle: subtitle,
	}

	insert := database.DB.NewInsert()
	insert.Model(&dbWeek)
	_, err := insert.Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
func (database *Database) GetWeekSubtitle(date time.Time, ctx context.Context) (string, error) {
	if date.IsZero() {
		return "", errors.New("date is required")
	}
	monday := getMonday(date)
	dbWeek := dbModels.WeekSubtitle{
		Date: monday,
	}
	query := database.DB.NewSelect()
	query.Model(&dbWeek)
	query.WherePK()
	err := query.Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			events, err := dataCollectors.DataCollectors.WeekGoogleCalenderAPI.GetEvents(monday, monday.AddDate(0, 0, 1))
			if err != nil {
				return "", err
			} else {
				// add if not exist cache
				for _, event := range events {
					fmt.Println(event.Summary)

					if matched := regexp.MustCompile(".*-Woche").MatchString(event.Summary); matched { //toconfig
						err = database.CreateWeekSubtitle(monday, event.Summary, ctx)
						if err != nil {
							return event.Summary, err
						}
						return event.Summary, nil
					}
				}
			}
		}
		return "", err
	}

	return dbWeek.Subtitle, nil
}

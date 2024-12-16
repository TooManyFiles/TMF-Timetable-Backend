package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

func (database *Database) CreateHolyday(startDate time.Time, endDate time.Time, name string, ctx context.Context) error {
	if startDate.IsZero() {
		return errors.New("startDate is required")
	}
	if endDate.IsZero() {
		return errors.New("endDate is required")
	}
	dbWeek := dbModels.Holyday{
		StartDate: startDate,
		EndDate:   endDate,
		Name:      name,
	}

	insert := database.DB.NewInsert()
	insert.Model(&dbWeek)
	_, err := insert.Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
func (database *Database) FetchHolydays(ctx context.Context) error {
	data, err := dataCollectors.DataCollectors.UntisClient.GetHolidays()
	if err != nil {
		return err
	}

	holydays := make([]dbModels.Holyday, len(data))
	for i, t := range data {
		startDateStr := strconv.Itoa(t.StartDate)
		startDate, err := time.Parse("20060102", startDateStr)
		if err != nil {
			return fmt.Errorf("error parsing date: %w", err)
		}
		endDateStr := strconv.Itoa(t.EndDate)
		endDate, err := time.Parse("20060102", endDateStr)
		if err != nil {
			return fmt.Errorf("error parsing date: %w", err)
		}
		teacher := dbModels.Holyday{
			Name:      t.LongName,
			StartDate: startDate,
			EndDate:   endDate,
		}
		holydays[i] = teacher
	}
	query := database.DB.NewInsert()
	query.On("CONFLICT (id) DO UPDATE")
	_, err = query.Model(&holydays).Exec(ctx)
	return err

}
func (database *Database) GetHolyday(date time.Time, ctx context.Context) (gen.Holyday, error) {
	if date.IsZero() {
		return gen.Holyday{}, errors.New("date is required")
	}

	dbHolyday := dbModels.Holyday{}
	query := database.DB.NewSelect()
	query.Model(&dbHolyday)
	query.Where("? BETWEEN start_date AND end_date", date)
	err := query.Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbHolyday.ToGen(), err
		}
		return dbHolyday.ToGen(), err
	}

	return dbHolyday.ToGen(), nil
}
func (database *Database) GetHolydays(date time.Time, days int, ctx context.Context) ([]gen.Holyday, error) {
	if date.IsZero() {
		return nil, errors.New("date is required")
	}

	// Calculate the end date for the date range
	endDate := date.AddDate(0, 0, days) // Adds 'days' to the date

	var dbHolydays []dbModels.Holyday
	err := database.DB.NewSelect().
		Model(&dbHolydays).
		Where("start_date <= ? AND end_date >= ?", endDate, date).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No holidays found
		}
		return nil, err // Return other errors
	}

	// Convert dbModels.Holyday to gen.Holyday
	holydays := make([]gen.Holyday, len(dbHolydays))
	for i, dbHolyday := range dbHolydays {
		holydays[i] = dbHolyday.ToGen()
	}

	return holydays, nil
}

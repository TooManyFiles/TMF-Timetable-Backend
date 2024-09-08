package db

import (
	"context"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

func (database *Database) CreateChoice(choice gen.Choice, ctx context.Context) (gen.Choice, error) {
	var dbChoice dbModels.Choice
	dbChoice.FromGen(choice)

	insert := database.DB.NewInsert()
	insert.Model(&dbChoice)
	_, err := insert.Exec(ctx)

	if err != nil {
		return gen.Choice{}, err
	}

	return dbChoice.ToGen(), nil
}

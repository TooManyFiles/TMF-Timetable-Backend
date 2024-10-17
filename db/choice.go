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
func (database *Database) CreateOrUpdateChoice(userId int, choiceId int, choice gen.Choice, ctx context.Context) (gen.Choice, error) {
	var dbChoice dbModels.Choice
	if choiceId == -1 {
		choice.Id = nil
	} else {
		choice.Id = &choiceId
	}
	choice.UserId = &userId
	dbChoice.FromGen(choice)
	insert := database.DB.NewInsert()
	insert.Model(&dbChoice)
	insert.On("CONFLICT (\"id\", \"userId\") DO UPDATE")
	insert.Where("\"choice\".\"userId\" = EXCLUDED.\"userId\"")
	_, err := insert.Exec(ctx)

	if err != nil {
		return gen.Choice{}, err
	}

	return dbChoice.ToGen(), nil
}
func (database *Database) GetChoiceByUserIdAndChoiceId(userId int, choiceId int, ctx context.Context) (gen.Choice, error) {
	dbChoice := dbModels.Choice{
		Id:     choiceId,
		UserId: userId,
	}
	query := database.DB.NewSelect()
	query.Model(&dbChoice)
	query.WherePK()
	err := query.Scan(ctx)

	if err != nil {
		return gen.Choice{}, err
	}

	return dbChoice.ToGen(), nil
}
func (database *Database) GetChoicesByUserId(userId int, ctx context.Context) ([]gen.Choice, error) {
	var dbChoices []dbModels.Choice

	err := database.DB.NewSelect().
		Model(&dbChoices).
		Where("\"userId\" = ?", userId).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	// Convert to gen.Choice slice
	genChoices := make([]gen.Choice, len(dbChoices))
	for i, dbChoice := range dbChoices {
		genChoices[i] = dbChoice.ToGen()
	}

	return genChoices, nil
}

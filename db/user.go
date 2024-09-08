package db

import (
	"context"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

func (database *Database) CreateUser(user gen.User, pwd string, ctx context.Context) (gen.User, error) {
	// unset DefaultChoice.Id to Prevent Collisions
	if user.DefaultChoice == nil {
		name := "Default"
		choice := map[string]interface{}{}
		user.DefaultChoice = &gen.Choice{
			Name:   &name,
			Choice: &choice,
		}
	} else {
		user.DefaultChoice.Id = nil
	}
	// Hash the password with bcrypt
	hashedPwd, err := hashPassword(pwd)
	if err != nil {
		return gen.User{}, err
	}
	var dbUser dbModels.User
	dbUser.FromGen(user)
	dbUser.PwdHash = hashedPwd

	insert := database.DB.NewInsert()
	insert.Model(&dbUser)
	_, err = insert.Exec(ctx)

	if err != nil {
		return gen.User{}, err
	}

	if user.DefaultChoice.Name == nil {
		name := "Default"
		user.DefaultChoice.Name = &name
	}
	if user.DefaultChoice.Choice == nil {
		choice := map[string]interface{}{}
		user.DefaultChoice.Choice = &choice
	}
	user.DefaultChoice.UserId = &dbUser.Id
	createdChoice, err := database.CreateChoice(*user.DefaultChoice, ctx)
	if err != nil {
		return dbUser.ToGen(), err
	}
	dbUser.DefaultChoiceId = *createdChoice.Id
	database.DB.NewUpdate().Model(&dbUser).WherePK().Exec(ctx)
	dbUser.DefaultChoice = &dbModels.Choice{}
	dbUser.DefaultChoice.FromGen(createdChoice)
	return dbUser.ToGen(), nil
}

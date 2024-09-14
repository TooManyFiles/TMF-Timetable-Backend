package db

import (
	"context"
	"database/sql"
	"errors"

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
func (database *Database) GetUserByID(id int, ctx context.Context) (gen.User, error) {
	var user dbModels.User
	query := database.DB.NewSelect()
	query.Model(&user)
	query.Where("\"user\".\"id\" = ?", id)
	query.Relation("DefaultChoice")
	err := query.Scan(ctx) //sql.ErrNoRows

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return gen.User{}, dbModels.ErrUserNotFound
		}
		return gen.User{}, err
	}
	return user.ToGen(), nil

}
func (database *Database) DeleteUserByID(id int, ctx context.Context) error {
	var user dbModels.User
	query := database.DB.NewDelete()
	query.Model(&user)
	query.Where("\"user\".\"id\" = ?", id)
	_, err := query.Exec(ctx) //sql.ErrNoRows
	if err != nil {
		return err
	}
	var choice dbModels.Choice
	query = database.DB.NewDelete()
	query.Model(&choice)
	query.Where("\"choice\".\"userId\" = ?", id)
	_, err = query.Exec(ctx) //sql.ErrNoRows
	if err != nil {
		return err
	}
	return nil

}
func (database *Database) GetUsers(ctx context.Context) ([]gen.User, error) {
	var users []dbModels.User
	query := database.DB.NewSelect()
	query.Model(&users)
	err := query.Scan(ctx) //sql.ErrNoRows
	if err != nil {
		return []gen.User{}, err
	}
	genUsers := make([]gen.User, len(users))
	for i, s := range users {
		genUsers[i] = s.ToGen()
	}
	return genUsers, err
}

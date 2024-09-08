package db

import (
	"context"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

func (database *Database) CreateUser(user gen.User, pwd string, ctx context.Context) (gen.User, error) {
	// Hash the password with bcrypt
	hashedPwd, err := hashPassword(pwd)
	if err != nil {
		return gen.User{}, err
	}

	var dbUser dbModels.User
	dbUser.FromGenUser(user)
	dbUser.PwdHash = hashedPwd

	insert := database.DB.NewInsert()
	insert.Model(&dbUser)
	_, err = insert.Exec(ctx)

	if err != nil {
		return gen.User{}, err
	}

	return dbUser.ToGenUser(), nil
}

package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

func checkUsernameCriteria(username string) bool {
	if len(username) >= 5 {
		return true
	}
	return false
}
func (database *Database) CreateUser(user gen.User, pwd string, ctx context.Context) (gen.User, error) {
	if !checkUsernameCriteria(user.Name) {
		return gen.User{}, dbModels.ErrUsernameNotMachRequirements
	}
	if len(pwd) == 0 {
		return gen.User{}, dbModels.ErrPasswordNotMachRequirements
	}
	// unset DefaultChoice.Id to Prevent Collisions
	if user.DefaultChoice == nil {
		name := "Default"
		choice := map[string]interface{}{}
		user.DefaultChoice = &gen.Choice{
			Name:   &name,
			Choice: &choice,
		}
	} else {
		// unset DefaultChoice.Id to Prevent Collisions
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
func (database *Database) fetchUser(user *dbModels.User, ctx context.Context) error {
	query := database.DB.NewSelect()
	query.Model(user)
	query.WherePK()
	query.Relation("DefaultChoice")
	err := query.Scan(ctx) //sql.ErrNoRows

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbModels.ErrUserNotFound
		}
		return err
	}
	return nil

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
func (database *Database) UpdateUntisLogin(genUser gen.User, untisName string, forename string, surname string, untisPWD string, key []byte, ctx context.Context) error {
	var user dbModels.User
	user.FromGen(genUser)
	user.UntisName = untisName
	err := dataCollectors.DataCollectors.UntisClient.SetupStudent(&user, forename, surname, untisPWD)
	if err != nil {
		return err
	}
	query := database.DB.NewUpdate()
	query.Model(&user)
	query.WherePK()
	query.Column("untis_pwd", "untis_name", "untis_id")

	encryptData, err := encrypt([]byte(untisPWD), key)
	if err != nil {
		return err
	}
	user.UntisPWD = base64.StdEncoding.EncodeToString(encryptData)
	_, err = query.Exec(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbModels.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (database *Database) GetUntisLogin(genUser gen.User, key []byte, ctx context.Context) (string, error) {
	var user dbModels.User
	user.FromGen(genUser)
	query := database.DB.NewSelect()
	query.Model(&user)
	query.WherePK()
	err := query.Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", dbModels.ErrUserNotFound
		}
		return "", err
	}

	untisPWD, err := base64.StdEncoding.DecodeString(user.UntisPWD)
	if err != nil {
		return "", err
	}

	decryptData, err := decrypt(untisPWD, key)
	if err != nil {
		return "", err
	}

	return string(decryptData), nil
}
func (database *Database) GetUntisLoginByCryptoKey(CryptoKey string, user gen.User, ctx context.Context) (string, error) {

	key, err := base64.StdEncoding.DecodeString(CryptoKey)
	if err != nil {
		return "", err
	}
	return database.GetUntisLogin(user, key, ctx)

}

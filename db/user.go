package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"

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

// Example UpdateUntisLogin using the new generic functions
func (database *Database) UpdateUntisLogin(genUser gen.User, untisName, forename, surname, untisPWD string, key []byte, ctx context.Context) error {
	var user dbModels.User
	user.FromGen(genUser)

	// Encrypt the password before storing
	encryptData, err := encrypt([]byte(untisPWD), key)
	if err != nil {
		return err
	}
	encodedPWD := base64.StdEncoding.EncodeToString(encryptData)

	// Call UntisClient setup
	untisId, personType, classId, err := dataCollectors.DataCollectors.UntisClient.SetupStudent(untisName, forename, surname, untisPWD)
	if err != nil {
		return err
	}
	// Update settings ind db
	if err := database.UpdateUserSetting(user.Id, "untis", "untisName", untisName, ctx); err != nil {
		return err
	}
	if err := database.UpdateUserSetting(user.Id, "untis", "untisPWD", encodedPWD, ctx); err != nil {
		return err
	}
	if err := database.UpdateUserSetting(user.Id, "untis", "userId", strconv.Itoa(untisId), ctx); err != nil {
		return err
	}
	if err := database.UpdateUserSetting(user.Id, "untis", "personType", strconv.Itoa(personType), ctx); err != nil {
		return err
	}
	if err := database.UpdateUserSetting(user.Id, "untis", "classId", strconv.Itoa(classId), ctx); err != nil {
		return err
	}
	return nil
}

func (database *Database) GetUntisLogin(genUser gen.User, key []byte, ctx context.Context) (string, string, error) {
	var user dbModels.User
	user.FromGen(genUser)

	// Retrieve the encrypted UntisPWD setting
	encPWD, err := database.GetUserSetting(user.Id, "untis", "untisPWD", ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", dbModels.ErrUserNotFound
		}
		return "", "", err
	}

	// Retrieve the UntisName setting
	untisName, err := database.GetUserSetting(user.Id, "untis", "untisName", ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", dbModels.ErrUserNotFound
		}
		return "", "", err
	}

	// Decode and decrypt the password
	untisPWD, err := base64.StdEncoding.DecodeString(encPWD)
	if err != nil {
		return "", "", err
	}
	if len(untisPWD) == 0 {
		return "", "", errors.New("no untis login")
	}
	decryptData, err := decrypt(untisPWD, key)
	if err != nil {
		return "", "", err
	}

	return untisName, string(decryptData), nil
}
func (database *Database) GetUntisLoginByCryptoKey(CryptoKey string, user gen.User, ctx context.Context) (string, string, error) {

	key, err := base64.StdEncoding.DecodeString(CryptoKey)
	if err != nil {
		return "", "", err
	}
	return database.GetUntisLogin(user, key, ctx)

}

func (database *Database) UpdateUser(user gen.User, ctx context.Context) error {
	var dbUser dbModels.User
	err := database.DB.NewSelect().Model(&dbUser).Where("id = ?", user.Id).Scan(ctx)
	if err != nil {
		return err
	}

	columns := []string{}

	if dbUser.Name != user.Name {
		dbUser.Name = user.Name
		columns = append(columns, "name")
	}
	if dbUser.Email != *user.Email {
		dbUser.Email = *user.Email
		columns = append(columns, "email")
	}
	if dbUser.Role != string(*user.Role) {
		dbUser.Role = string(*user.Role)
		columns = append(columns, "role")
	}
	if dbUser.DefaultChoiceId != *user.DefaultChoice.Id {
		dbUser.DefaultChoiceId = *user.DefaultChoice.Id
		columns = append(columns, "defaultChoice")
	}
	if len(columns) > 0 {
		_, err := database.DB.NewUpdate().
			Model(&dbUser).
			Column(columns...).
			WherePK().
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

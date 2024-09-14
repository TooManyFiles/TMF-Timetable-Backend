package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func (database *Database) GetTeachers(ctx context.Context) ([]gen.Teacher, error) {
	query := database.DB.NewSelect()
	teachers := make([]dbModels.Teacher, 0)

	query.Model(&teachers)
	err := query.Scan(ctx)
	genTeachers := make([]gen.Teacher, len(teachers))
	for i, t := range teachers {
		genTeachers[i] = t.ToGen()
	}
	return genTeachers, err
}
func (database *Database) GetSubjects(ctx context.Context) ([]gen.Subject, error) {
	query := database.DB.NewSelect()
	subjects := make([]dbModels.Subject, 0)

	query.Model(&subjects)
	err := query.Scan(ctx)
	genSubjects := make([]gen.Subject, len(subjects))
	for i, s := range subjects {
		genSubjects[i] = s.ToGen()
	}
	return genSubjects, err
}
func (database *Database) GetRooms(ctx context.Context) ([]gen.Room, error) {
	query := database.DB.NewSelect()
	rooms := make([]dbModels.Room, 0)

	query.Model(&rooms)
	err := query.Scan(ctx)
	genRooms := make([]gen.Room, len(rooms))
	for i, r := range rooms {
		genRooms[i] = r.ToGen()
	}
	return genRooms, err
}
func (database *Database) GetClasses(ctx context.Context) ([]gen.Class, error) {
	query := database.DB.NewSelect()
	classes := make([]dbModels.Class, 0)

	query.Model(&classes)
	err := query.Scan(ctx)
	genClasses := make([]gen.Class, len(classes))
	for i, c := range classes {
		genClasses[i] = c.ToGen()
	}
	return genClasses, err
}
func (database *Database) GetLesson(filter dbModels.LessonFilter, ctx context.Context) ([]gen.Lesson, error) {
	if filter.User.Id == 0 {
		return nil, errors.New("user ID is required to fetch lessons")
	}
	var choice dbModels.Choice
	// get Choice
	if filter.Choice.Id == 0 {
		if filter.Choice.Choice == "" {
			query := database.DB.NewSelect()
			query.Model(&filter.User)
			query.WherePK()
			query.Relation("DefaultChoice")
			err := query.Scan(ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, dbModels.ErrUserNotFound
				}
				return nil, err
			}
			choice = *filter.User.DefaultChoice
		} else {
			choice = filter.Choice
			userQuery := database.DB.NewSelect()
			userQuery.Model(&filter.User)
			userQuery.WherePK()
			err := userQuery.Scan(ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, dbModels.ErrUserNotFound
				}
				return nil, err
			}
			//TODO: check class
		}

	} else if filter.Choice.Id != 0 {
		userQuery := database.DB.NewSelect()
		userQuery.Model(&filter.User)
		userQuery.WherePK()
		err := userQuery.Scan(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, dbModels.ErrUserNotFound
			}
			return nil, err
		}

		query := database.DB.NewSelect()
		query.Model(&filter.Choice)
		query.WherePK()
		if filter.User.Role != string(gen.UserRoleAdmin) {
			query.Where("\"choice\".\"userId\" = ?", filter.User.Id) //TODO: disable if Admin
		}
		err = query.Scan(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, dbModels.ErrChoiceNotFound
			}
			return nil, err
		}
		choice = filter.Choice
	}

	lessonQuery := database.DB.NewSelect()
	lessons := make([]dbModels.Lesson, 0)
	lessonQuery.Model(&lessons)

	var result map[string]interface{}
	parsingError := json.Unmarshal([]byte(choice.Choice), &result)

	if parsingError != nil || choice.Choice == "" || len(result) == 0 {
		lessonQuery.Where("\"lesson\".\"classes\" @> ?", pgdialect.Array(filter.User.Classes))
	} else {

		for key, value := range result {
			if classID, err := strconv.Atoi(key); err == nil {
				if classID > 0 {
					lessonQuery.WhereOr("(\"lesson\".\"classes\" = ? AND \"lesson\".\"subjects\" @> ?)", classID, pgdialect.Array(value))
				} else { //TODO: If a Class ID is present as a negative as well as a positive value only the positive should be used.
					lessonQuery.WhereOr("(\"lesson\".\"classes\" = ? AND NOT \"lesson\".\"subjects\" @> ?)", classID, pgdialect.Array(value))
				}
			}
		}
	}
	err := lessonQuery.Scan(ctx)
	genLesson := make([]gen.Lesson, len(lessons))
	for i, c := range lessons {
		genLesson[i] = c.ToGen()
	}
	return genLesson, err
}

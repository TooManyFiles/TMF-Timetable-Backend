package db

import (
	"context"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
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

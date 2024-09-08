package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DatabaseConfig struct {
	PG pgdriver.Option
}
type Database struct {
	DatabaseConfig
	DB *bun.DB
}

func NewDatabase(config DatabaseConfig) Database {
	database := Database{DatabaseConfig: config}
	ctx := context.Background()

	// Open a PostgreSQL database.
	pgdb := sql.OpenDB(pgdriver.NewConnector(database.PG))

	// Create a Bun db on top of it.
	database.DB = bun.NewDB(pgdb, pgdialect.New())

	// Print all queries to stdout.
	database.DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	var rnd float64
	log.Println("Connected..")
	if err := database.DB.NewSelect().ColumnExpr("random()").Scan(ctx, &rnd); err != nil {
		panic(err)
	}

	err := database.createSchema(ctx)
	if err != nil {
		panic(err)
	}

	return database
}

func (database *Database) createSchema(ctx context.Context) error {
	// List of models to create
	models := []interface{}{
		&Class{},
		&User{},
		&Teacher{},
		&Lesson{},
		&Room{},
		&Subject{},
		&Choice{},
		&Menu{},
	}

	for _, model := range models {
		// Create the table if it does not exist
		table := database.DB.NewCreateTable()
		table.IfNotExists()
		table.Model(model)
		res, err := table.Exec(ctx)
		if err != nil {
			return err
		}
		log.Println(res)
	}
	log.Println("Schema created.")
	return nil
}

// FetchMenuForDate fetches a menu from the database based on the given date
func (database *Database) FetchMenuForDate(startDate time.Time, days int, ctx context.Context) ([]gen.Menu, error) {

	// Calculate the end date by adding 'days - 1' to startDate
	endDate := startDate.AddDate(0, 0, days-1)

	var dbMenus []Menu
	err := database.DB.NewSelect().
		Model(&dbMenus).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(ctx)
	if err != nil {
		log.Println("Error fetching menus:", err)
		return []gen.Menu{}, err
	}
	if dbMenus == nil {
		return []gen.Menu{}, nil
	}

	// Convert db.Menu to gen.Menu in one pass
	menus := make([]gen.Menu, len(dbMenus))
	for i, dbMenu := range dbMenus {
		menus[i] = gen.Menu{
			Cookteam:    &dbMenu.Cookteam,
			Date:        &openapi_types.Date{Time: dbMenu.Date},
			Dessert:     &dbMenu.Dessert,
			Garnish:     &dbMenu.Garnish,
			MainDish:    &dbMenu.MainDish,
			MainDishVeg: &dbMenu.MainDishVeg,
		}
	}
	return menus, nil
}

package db

import (
	"context"
	"database/sql"
	"log"

	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
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
		&dbModels.Class{},
		&dbModels.User{},
		&dbModels.Teacher{},
		&dbModels.Lesson{},
		&dbModels.Room{},
		&dbModels.Subject{},
		&dbModels.Choice{},
		&dbModels.Menu{},
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

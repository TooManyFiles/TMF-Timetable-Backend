package db

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// FetchMenuForDate fetches a menu from the database based on the given date
func (database *Database) FetchMenuForDate(startDate time.Time, days int, ctx context.Context) ([]gen.Menu, error) {

	// Calculate the end date by adding 'days - 1' to startDate
	endDate := startDate.AddDate(0, 0, days-1)

	var dbMenus []dbModels.Menu
	err := database.DB.NewSelect().
		Model(&dbMenus).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(ctx)
	if err != nil {
		log.Println("Error fetching menus:", err)
		return []gen.Menu{}, err
	}
	if dbMenus == nil {
		menus, err := dataCollectors.DataCollectors.TFfoodplanAPI.GetForRange(startDate, days)
		if err != nil {
			log.Println("Error fetching menus:", err)
			return []gen.Menu{}, err
		}
		_, err = database.DB.NewInsert().Model(&menus).Exec(ctx)
		if err != nil {
			log.Println("Error inserting new menus:", err)
			return []gen.Menu{}, err
		}

		dbMenus = getFirstNMenus(menus, days)
	} else if len(dbMenus) < days {
		log.Println(len(dbMenus), days)
		// Fetch new menus from the API
		menus, err := dataCollectors.DataCollectors.TFfoodplanAPI.GetForRange(startDate, days)
		if err != nil {
			log.Println("Error fetching menus from API:", err)
			return []gen.Menu{}, err
		}
		log.Println(menus)
		// Upsert (update or insert) fetched menus into the database
		_, err = database.DB.NewInsert().
			Model(&menus).
			On("CONFLICT (date) DO UPDATE").
			Set("cookteam = EXCLUDED.cookteam").
			Set("main_dish = EXCLUDED.main_dish").
			Set("main_dish_veg = EXCLUDED.main_dish_veg").
			Set("garnish = EXCLUDED.garnish").
			Set("dessert = EXCLUDED.dessert").
			Exec(ctx)
		if err != nil {
			log.Println("Error upserting menus:", err)
			return []gen.Menu{}, err
		}
		dbMenus = getFirstNMenus(menus, days)
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

func getFirstNMenus(menus []dbModels.Menu, n int) []dbModels.Menu {
	// Sort menus by date
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].Date.Before(menus[j].Date)
	})

	// Ensure n is within the bounds of the slice length
	if n > len(menus) {
		n = len(menus)
	}

	// Return the first n elements
	return menus[:n]
}

//go:build tools
// +build tools

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=tools/codegen-config/models.yml TMF-Timetable-Docs/swagger.yml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=tools/codegen-config/server.yml TMF-Timetable-Docs/swagger.yml
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api"
	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/config"
	tffoodplanapi "github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/TFfoodplanAPI"
	"github.com/TooManyFiles/TMF-Timetable-Backend/db"
)

var database db.Database

func main() {
	database = db.NewDatabase(config.DatabaseConfig)
	menu()

	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := api.NewServer(database)

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := gen.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	// // And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
func menu() {
	api := &tffoodplanapi.TFfoodplanAPI{
		URL: "http://www.treffpunkt-fanny.de/images/stories/dokumente/Essensplaene/api/TFfoodplanAPI.php",
	}

	// Fetch menu for a specific date
	menu, err := api.GetForDate(time.Date(2024, 9, 10, 0, 0, 0, 0, time.Local))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Fetched menu:", menu)
	ctx := context.Background()

	database.DB.NewInsert().Model(&menu).Exec(ctx)

	//

	// Update menu (if needed)
	updatedMenu, err := api.Update(menu)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated menu:", updatedMenu)
	startDate := time.Now()
	dataCount := 3 // Fetch menus for 3 days
	menus, err := api.GetForRange(startDate, dataCount)
	if err != nil {
		log.Fatal("Failed to fetch menus:", err)
	}

	for _, menu := range menus {
		fmt.Printf("Date: %s, Main Dish: %s\n", menu.Date.Format("2.1.2006"), menu.MainDish)
	}
}

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
	"github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors"
	"github.com/TooManyFiles/TMF-Timetable-Backend/db"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	"github.com/rs/cors"
)

var database db.Database

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	dataCollectors.InitDataCollectors()
	initDB()
	initServer()

}

func initDB() {
	database = db.NewDatabase(config.Config.DatabaseConfig)
	querry := database.DB.NewInsert()
	less := dbModels.Lesson{
		Subjects: []string{"1", "2", "4", " 5"},
	}
	querry.Model(&less)
	_, err := querry.Exec(context.Background())
	if err != nil {
		log.Println(err.Error())

	}
}

func initServer() {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := api.NewServer(database)

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := gen.HandlerFromMux(server, r)
	handler := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions, http.MethodPut},
		Logger:         log.Default(),
		Debug:          true,
		AllowedHeaders: []string{"*"},
	}).Handler(h)

	s := &http.Server{
		Handler:                      handler,
		Addr:                         "0.0.0.0:8080",
		DisableGeneralOptionsHandler: true,
	}

	// // And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}

func menu() {
	// Fetch menu for a specific date
	menu, err := dataCollectors.DataCollectors.TFfoodplanAPI.GetForDate(time.Date(2024, 9, 10, 0, 0, 0, 0, time.Local))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Fetched menu:", menu)
	ctx := context.Background()

	database.DB.NewInsert().Model(&menu).Exec(ctx)

	//

	// Update menu (if needed)
	updatedMenu, err := dataCollectors.DataCollectors.TFfoodplanAPI.Update(menu)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated menu:", updatedMenu)
	startDate := time.Now()
	dataCount := 3 // Fetch menus for 3 days
	menus, err := dataCollectors.DataCollectors.TFfoodplanAPI.GetForRange(startDate, dataCount)
	if err != nil {
		log.Fatal("Failed to fetch menus:", err)
	}

	for _, menu := range menus {
		fmt.Printf("Date: %s, Main Dish: %s\n", menu.Date.Format("2.1.2006"), menu.MainDish)
	}
}

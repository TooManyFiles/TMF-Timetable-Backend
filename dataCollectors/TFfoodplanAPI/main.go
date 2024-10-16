package tffoodplanapi

import (
	"fmt"
	"log"
	"net/url"
	"time"

	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	"github.com/go-resty/resty/v2" // Import Resty
)

type TFfoodplanAPI struct {
	URL string
}

// getForDate fetches the menu for a specific date from the external API
func (api *TFfoodplanAPI) GetForDate(date time.Time) (dbModels.Menu, error) {
	client := resty.New()

	// Format the date for the API
	dateStr := date.Format("02.01.2006")

	// Construct the URL with query parameters
	apiUrl := fmt.Sprintf("%s?dateFormat=j/n/Y&dataMode=days&dataCount=1&dataFromTime=%s", api.URL, url.QueryEscape(dateStr))
	log.Println(apiUrl)
	// Define the structure to hold the API response
	var apiResponse []struct {
		Date        string `json:"date"`
		Cookteam    string `json:"cookteam"`
		MainDish    string `json:"mainDish"`
		MainDishVeg string `json:"mainDishVeg"`
		Garnish     string `json:"garnish"`
		Dessert     string `json:"dessert"`
	}

	// Make the request and automatically unmarshal the response JSON into the struct
	_, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&apiResponse). // Automatically unmarshal the JSON into apiResponse
		Get(apiUrl)

	if err != nil {
		log.Println("Error fetching data:", err)
		return dbModels.Menu{}, err
	}

	// Check if the response is empty
	if len(apiResponse) == 0 {
		log.Println("No data returned from API")
		return dbModels.Menu{}, fmt.Errorf("no data returned from API")
	}
	// Check if the response is empty
	if len(apiResponse) == 0 {
		log.Println("No data returned from API for date:", date)
		// Return a menu indicating the date was skipped
		return dbModels.Menu{
			Date:         date,
			NotAPIServed: true,
		}, nil
	}
	// Convert the API response to the db.Menu struct
	menuDate, _ := time.Parse("2/1/2006", apiResponse[0].Date) // Parse the date
	menu := dbModels.Menu{
		Date:        menuDate,
		Cookteam:    apiResponse[0].Cookteam,
		MainDish:    apiResponse[0].MainDish,
		MainDishVeg: apiResponse[0].MainDishVeg,
		Garnish:     apiResponse[0].Garnish,
		Dessert:     apiResponse[0].Dessert,
	}

	return menu, nil
}

// Update sends the updated menu to the API
func (api *TFfoodplanAPI) Update(menu dbModels.Menu) (dbModels.Menu, error) {
	return api.GetForDate(menu.Date)
}

// getForRange fetches menus for a range of dates from the external API
func (api *TFfoodplanAPI) GetForRange(startDate time.Time, dataCount int) ([]dbModels.Menu, error) {
	client := resty.New()

	// Format the date for the API
	startDateStr := startDate.Format("2.1.2006")

	// Construct the URL with query parameters for a date range
	apiUrl := fmt.Sprintf("%s?dateFormat=j/n/Y&dataMode=days&dataCount=%d&dataFromTime=%s", api.URL, dataCount, url.QueryEscape(startDateStr))

	// Define the structure to hold the API response
	var apiResponse []struct {
		Date        string `json:"date"`
		Cookteam    string `json:"cookteam"`
		MainDish    string `json:"mainDish"`
		MainDishVeg string `json:"mainDishVeg"`
		Garnish     string `json:"garnish"`
		Dessert     string `json:"dessert"`
	}

	// Make the request and automatically unmarshal the response JSON into the struct
	_, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&apiResponse). // Automatically unmarshal the JSON into apiResponse
		Get(apiUrl)

	if err != nil {
		log.Println("Error fetching data:", err)
		return nil, err
	}

	// Initialize a slice to hold the menus
	menus := make([]dbModels.Menu, 0, dataCount)

	// Create a map for quick lookup of dates received from API
	receivedDates := make(map[time.Time]bool)

	// Iterate over the API response using range and convert each entry to db.Menu
	for _, apiMenu := range apiResponse {
		menuDate, _ := time.Parse("2/1/2006", apiMenu.Date) // Parse the date
		menu := dbModels.Menu{
			Date:        menuDate,
			Cookteam:    apiMenu.Cookteam,
			MainDish:    apiMenu.MainDish,
			MainDishVeg: apiMenu.MainDishVeg,
			Garnish:     apiMenu.Garnish,
			Dessert:     apiMenu.Dessert,
		}
		menus = append(menus, menu)
		receivedDates[menuDate] = true
	}
	startDate = startDate.UTC().Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, dataCount)
	// Identify and create skipped dates
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if !receivedDates[currentDate] {
			menus = append(menus, dbModels.Menu{
				Date:         currentDate,
				NotAPIServed: true,
			})
		}
	}
	return menus, nil
}

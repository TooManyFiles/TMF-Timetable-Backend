package dataCollectors

import (
	"github.com/TooManyFiles/TMF-Timetable-Backend/config"
	tffoodplanapi "github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/TFfoodplanAPI"
	googleapi "github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/googleAPI"
	untisDataCollectors "github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/untisDataCollectors"
)

var DataCollectors DataCollectorsStruct

type DataCollectorsStruct struct {
	TFfoodplanAPI         tffoodplanapi.TFfoodplanAPI
	UntisClient           untisDataCollectors.UntisClient
	WeekGoogleCalenderAPI googleapi.GoogleCalenderAPI
}

func InitDataCollectors() {
	DataCollectors.TFfoodplanAPI = tffoodplanapi.TFfoodplanAPI{
		URL: config.Config.DataCollectors.TFfoodplanAPIURL,
	}
	DataCollectors.UntisClient, _ = untisDataCollectors.Init(config.Config.DataCollectors.UntisApiConfig) //TODO: error handling
	DataCollectors.WeekGoogleCalenderAPI = googleapi.GoogleCalenderAPI{
		ApiKey:     config.Config.DataCollectors.WeekGoogleCalenderAPIConfig.ApiKey,
		CalendarID: config.Config.DataCollectors.WeekGoogleCalenderAPIConfig.CalendarID,
	}
}

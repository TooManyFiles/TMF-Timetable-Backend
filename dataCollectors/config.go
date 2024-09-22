package dataCollectors

import (
	"github.com/TooManyFiles/TMF-Timetable-Backend/config"
	tffoodplanapi "github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/TFfoodplanAPI"
	untisDataCollectors "github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/untisDataCollectors"
)

var DataCollectors DataCollectorsStruct

type DataCollectorsStruct struct {
	TFfoodplanAPI tffoodplanapi.TFfoodplanAPI
	UntisClient   untisDataCollectors.UntisClient
}

func InitDataCollectors() {
	DataCollectors.TFfoodplanAPI = tffoodplanapi.TFfoodplanAPI{
		URL: config.Config.DataCollectors.TFfoodplanAPIURL,
	}
	DataCollectors.UntisClient, _ = untisDataCollectors.Init(config.Config.DataCollectors.UntisApiConfig)
}

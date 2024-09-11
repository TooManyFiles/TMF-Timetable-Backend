package dataCollectors

import (
	"github.com/TooManyFiles/TMF-Timetable-Backend/config"
	tffoodplanapi "github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/TFfoodplanAPI"
)

var DataCollectors DataCollectorsStruct

type DataCollectorsStruct struct {
	TFfoodplanAPI tffoodplanapi.TFfoodplanAPI
}

func InitDataCollectors() {
	DataCollectors.TFfoodplanAPI = tffoodplanapi.TFfoodplanAPI{
		URL: config.Config.DataCollectors.TFfoodplanAPIURL,
	}
}

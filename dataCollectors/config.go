package dataCollectors

import tffoodplanapi "github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors/TFfoodplanAPI"

var DataCollectors DataCollectorsStruct

type DataCollectorsStruct struct {
	TFfoodplanAPI tffoodplanapi.TFfoodplanAPI
}

func InitDataCollectors() {
	DataCollectors.TFfoodplanAPI = tffoodplanapi.TFfoodplanAPI{
		URL: "http://www.treffpunkt-fanny.de/images/stories/dokumente/Essensplaene/api/TFfoodplanAPI.php",
	}
}

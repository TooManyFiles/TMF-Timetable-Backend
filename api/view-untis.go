package api

import (
	"context"
	"fmt"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/db"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

type UntisProviderSettings struct {
	// Choice Choice of subjects for the classes. {class:[subjects]}
	// - If a class has a empty array as a choice all subjects should be shown.
	// - If the Class ID is negative it the the choice is a blacklist.
	// - If a Class ID is present as a negative as well as a positive value only the positive should be used.
	Choice *gen.Choice `json:"Choice,omitempty"`
}

func (server Server) UntisView(user gen.User, claims *db.Claims, providerSettings UntisProviderSettings, startdate time.Time, enddate time.Time, fetchLesson bool, ctx context.Context) ([]gen.Lesson, error) {
	_, untis_pwd, err := server.DB.GetUntisLoginByCryptoKey(claims.CryptoKey, user, ctx)
	if err != nil {
		return nil, err
	}
	if fetchLesson {
		for _, classId := range *user.Classes {
			err = server.DB.FetchLesson(user, untis_pwd, classId, startdate, enddate, ctx)
			if err != nil {
				fmt.Println("Failed to FetchLesson: " + err.Error())
			}
		}
	}
	lessonFilter := dbModels.LessonFilter{
		User:      (&dbModels.User{}).FromGen(user),
		StartDate: startdate,
		EndDate:   enddate,
	}
	if providerSettings.Choice != nil {
		lessonFilter.Choice = (&dbModels.Choice{}).FromGen(*providerSettings.Choice)
	}
	resp, err := server.DB.GetLesson(lessonFilter, ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

package untisDataCollectors

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/Mr-Comand/goUntisAPI/structs"
	"github.com/Mr-Comand/goUntisAPI/untisApi"
	"github.com/TooManyFiles/TMF-Timetable-Backend/config"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

type UntisClient struct {
	staticClient  *untisApi.Client
	dynamicClient *untisApi.Client
}

func Init(apiConfig structs.ApiConfig) (UntisClient, error) {
	untisClient := UntisClient{
		staticClient:  untisApi.NewClient(apiConfig, log.Default(), config.Config.DataCollectors.Logging.UntisApi.StaticClient),
		dynamicClient: untisApi.NewClient(apiConfig, log.Default(), config.Config.DataCollectors.Logging.UntisApi.DynamicClient),
	}
	err := untisClient.staticClient.Authenticate()
	if err != nil {
		return UntisClient{}, err
	}
	untisClient.staticClient.Test()
	return untisClient, nil
}

func (untisClient UntisClient) reAuthenticate() error {
	err := untisClient.staticClient.Test()
	if err != nil {
		var rpcerr *structs.RPCError
		if errors.As(err, &rpcerr) && rpcerr.Code == -8520 {
			return untisClient.staticClient.Authenticate()
		}
		return err
	}
	return nil
}
func (untisClient UntisClient) GetTeachers() ([]structs.Teacher, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	teachers, err := untisClient.staticClient.GetTeachers()
	if err != nil {
		return nil, err
	}
	return teachers, nil
}
func (untisClient UntisClient) GetSubjects() ([]structs.Subject, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	subjects, err := untisClient.staticClient.GetSubjects()
	if err != nil {
		return nil, err
	}
	return subjects, nil
}
func (untisClient UntisClient) GetRooms() ([]structs.Room, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	subjects, err := untisClient.staticClient.GetRooms()
	if err != nil {
		return nil, err
	}
	return subjects, nil
}
func (untisClient UntisClient) GetClasses() ([]structs.Class, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	classes, err := untisClient.staticClient.GetClasses()
	if err != nil {
		return nil, err
	}
	return classes, nil
}
func (untisClient UntisClient) GetLessonsByClass(class dbModels.Class, startDate time.Time, endDate time.Time) ([]structs.Period, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	body := structs.GetTimetableRequest{
		Element: structs.GetTimetableRequestElement{
			Type: 1,
			Id:   class.Id,
		},
		ShowBooking:   true,
		ShowInfo:      true,
		ShowLsText:    true,
		ShowSubstText: true,
		ShowLsNumber:  true,
	}
	if !startDate.IsZero() {
		body.StartDate, _ = strconv.Atoi(startDate.Local().Format("20060102"))
	}
	if endDate.IsZero() {
		body.EndDate, _ = strconv.Atoi(startDate.AddDate(0, 0, 7).Local().Format("20060102"))
	} else {
		body.StartDate, _ = strconv.Atoi(startDate.Local().Format("20060102"))
	}
	lessons, err := untisClient.staticClient.GetTimetable(body)
	if err != nil {
		return nil, err
	}
	return lessons, nil
}

func (untisClient UntisClient) GetLessonsByStudent(student dbModels.User, untisPWD string, startDate time.Time, endDate time.Time) ([]structs.Period, error) {
	dynamicClient := untisApi.NewClient(untisClient.dynamicClient.ApiConfig, log.Default(), config.Config.DataCollectors.Logging.UntisApi.DynamicClient)
	dynamicClient.ApiConfig.User = student.UntisName
	dynamicClient.ApiConfig.Password = untisPWD
	err := dynamicClient.Authenticate()
	if err != nil {
		dynamicClient.Logout()
		return nil, err
	}
	body := structs.GetTimetableRequest{
		Element: structs.GetTimetableRequestElement{
			Type: 5,
			Id:   student.UntisId,
		},
		ShowBooking:   true,
		ShowInfo:      true,
		ShowLsText:    true,
		ShowSubstText: true,
		ShowLsNumber:  true,
	}
	if !startDate.IsZero() {
		body.StartDate, _ = strconv.Atoi(startDate.Local().Format("20060102"))
	}
	if endDate.IsZero() {
		body.EndDate, _ = strconv.Atoi(startDate.AddDate(0, 0, 7).Local().Format("20060102"))
	} else {
		body.StartDate, _ = strconv.Atoi(startDate.Local().Format("20060102"))
	}
	lessons, err := dynamicClient.GetTimetable(body)
	if err != nil {
		dynamicClient.Logout()
		return nil, err
	}
	dynamicClient.Logout()
	return lessons, nil
}

// Function to search for a person by foreName and longName
func findPerson(students []structs.Student, foreName string, longName string) *structs.Student {
	for _, person := range students {
		if person.ForeName == foreName && person.LongName == longName {
			return &person
		}
	}
	return nil // Return nil if not found
}

var ErrStudentNotFound = errors.New("student not found")

func (untisClient UntisClient) SetupStudent(user *dbModels.User, forename string, surname string, untisPWD string) error {
	err := untisClient.reAuthenticate()
	if err != nil {
		return err
	}
	students, err := untisClient.staticClient.GetStudents()
	if err != nil {
		return err
	}
	student := findPerson(students, forename, surname)
	dynamicClient := untisApi.NewClient(untisClient.dynamicClient.ApiConfig, log.Default(), config.Config.DataCollectors.Logging.UntisApi.DynamicClient)
	dynamicClient.ApiConfig.User = user.UntisName
	dynamicClient.ApiConfig.Password = untisPWD
	err = dynamicClient.Authenticate()
	if err != nil {
		dynamicClient.Logout()
		return err
	}
	if student.ID == dynamicClient.PersonID {
		user.UntisId = dynamicClient.PersonID
		dynamicClient.Logout()
		return nil
	} else {
		dynamicClient.Logout()
		return ErrStudentNotFound
	}
}

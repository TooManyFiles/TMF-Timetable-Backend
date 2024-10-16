package untisDataCollectors

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/Mr-Comand/goUntisAPI/structs"
	"github.com/Mr-Comand/goUntisAPI/untisApi"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
)

type UntisClient struct {
	staticClient  *untisApi.Client
	dynamicClient *untisApi.Client
}

func Init(apiConfig structs.ApiConfig) (UntisClient, error) {
	untisClient := UntisClient{
		staticClient:  untisApi.NewClient(apiConfig, log.Default(), untisApi.DEBUG, true),
		dynamicClient: untisApi.NewClient(apiConfig, log.Default(), untisApi.DEBUG, true),
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

func (untisClient UntisClient) GetLessonsByStudent(UntisName string, untisPWD string, startDate time.Time, endDate time.Time, classId int) ([]structs.Period, error) {
	dynamicClient := untisApi.NewClient(untisClient.dynamicClient.ApiConfig, log.Default(), untisApi.DEBUG, true)
	dynamicClient.ApiConfig.User = UntisName
	dynamicClient.ApiConfig.Password = untisPWD
	err := dynamicClient.Authenticate()
	if err != nil {
		dynamicClient.Logout()
		return nil, err
	}
	body := structs.GetTimetableRequest{
		Element: structs.GetTimetableRequestElement{
			Type: 1,
			Id:   classId,
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
		body.EndDate, _ = strconv.Atoi(endDate.Local().Format("20060102"))
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

func (untisClient UntisClient) SetupStudent(untisName, forename, surname, untisPWD string) (int, int, int, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return 0, 0, 0, err
	}
	students, err := untisClient.staticClient.GetStudents()
	if err != nil {
		return 0, 0, 0, err
	}
	student := findPerson(students, forename, surname)
	if student == nil {
		return 0, 0, 0, ErrStudentNotFound
	}
	dynamicClient := untisApi.NewClient(untisClient.dynamicClient.ApiConfig, log.Default(), untisApi.DEBUG, true)
	dynamicClient.ApiConfig.User = untisName
	dynamicClient.ApiConfig.Password = untisPWD
	err = dynamicClient.Authenticate()
	if err != nil {
		dynamicClient.Logout()
		return 0, 0, 0, err
	}

	if student.ID == dynamicClient.PersonID {
		personType := dynamicClient.PersonType
		personID := dynamicClient.PersonID
		klasseID := dynamicClient.KlasseId
		dynamicClient.Logout()
		return personID, personType, klasseID, nil
	} else {
		dynamicClient.Logout()
		return 0, 0, 0, ErrStudentNotFound
	}
}

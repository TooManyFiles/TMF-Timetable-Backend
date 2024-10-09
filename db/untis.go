package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/TooManyFiles/TMF-Timetable-Backend/dataCollectors"
	dbModels "github.com/TooManyFiles/TMF-Timetable-Backend/db/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func (database *Database) FetchTeachers(ctx context.Context) error {
	data, err := dataCollectors.DataCollectors.UntisClient.GetTeachers()
	if err != nil {
		return err
	}
	teachers := make([]dbModels.Teacher, len(data))
	for i, t := range data {
		teacher := dbModels.Teacher{
			Id:        t.ID,
			Name:      t.LongName,
			FirstName: t.ForeName,
			ShortName: t.Name,
		}
		teachers[i] = teacher
	}
	querya := database.DB.NewInsert()
	_, err = querya.Model(&teachers).Exec(ctx)
	return err
}
func (database *Database) GetTeachers(ctx context.Context) ([]gen.Teacher, error) {
	query := database.DB.NewSelect()
	teachers := make([]dbModels.Teacher, 0)

	query.Model(&teachers)
	err := query.Scan(ctx)
	genTeachers := make([]gen.Teacher, len(teachers))
	for i, t := range teachers {
		genTeachers[i] = t.ToGen()
	}
	return genTeachers, err
}
func (database *Database) FetchSubjects(ctx context.Context) error {
	data, err := dataCollectors.DataCollectors.UntisClient.GetSubjects()
	if err != nil {
		return err
	}
	subjects := make([]dbModels.Subject, len(data))
	for i, t := range data {
		subject := dbModels.Subject{
			Id:        t.ID,
			Name:      t.LongName,
			ShortName: t.Name,
		}
		subjects[i] = subject
	}
	querya := database.DB.NewInsert()
	_, err = querya.Model(&subjects).Exec(ctx)
	return err
}
func (database *Database) GetSubjects(ctx context.Context) ([]gen.Subject, error) {
	query := database.DB.NewSelect()
	subjects := make([]dbModels.Subject, 0)

	query.Model(&subjects)
	err := query.Scan(ctx)
	genSubjects := make([]gen.Subject, len(subjects))
	for i, s := range subjects {
		genSubjects[i] = s.ToGen()
	}
	return genSubjects, err
}
func (database *Database) GetRooms(ctx context.Context) ([]gen.Room, error) {
	query := database.DB.NewSelect()
	rooms := make([]dbModels.Room, 0)

	query.Model(&rooms)
	err := query.Scan(ctx)
	genRooms := make([]gen.Room, len(rooms))
	for i, r := range rooms {
		genRooms[i] = r.ToGen()
	}
	return genRooms, err
}
func (database *Database) FetchRooms(ctx context.Context) error {
	data, err := dataCollectors.DataCollectors.UntisClient.GetRooms()
	if err != nil {
		return err
	}
	rooms := make([]dbModels.Room, len(data))
	for i, t := range data {
		room := dbModels.Room{
			Id:                    t.ID,
			Name:                  t.Name,
			AdditionalInformation: t.LongName,
		}
		rooms[i] = room
	}
	querya := database.DB.NewInsert()
	_, err = querya.Model(&rooms).Exec(ctx)
	return err
}
func (database *Database) GetClasses(ctx context.Context) ([]gen.Class, error) {
	query := database.DB.NewSelect()
	classes := make([]dbModels.Class, 0)

	query.Model(&classes)
	err := query.Scan(ctx)
	genClasses := make([]gen.Class, len(classes))
	for i, c := range classes {
		genClasses[i] = c.ToGen()
	}
	return genClasses, err
}
func (database *Database) FetchClasses(ctx context.Context) error {
	data, err := dataCollectors.DataCollectors.UntisClient.GetClasses()
	if err != nil {
		return err
	}
	rooms := make([]dbModels.Class, len(data))
	for i, t := range data {
		teacher := dbModels.Class{
			Id:                 t.ID,
			Name:               t.Name,
			MainTeacherId:      t.Teacher1,
			SecondaryTeacherId: t.Teacher2,
		}
		rooms[i] = teacher
	}
	querya := database.DB.NewInsert()
	_, err = querya.Model(&rooms).Exec(ctx)
	return err
}

// MergeDateAndTime takes a date in YYYYMMDD format and a start time in HHMM format
// and returns a time.Time object.
func MergeDateAndTime(periodDate int, periodTime int) (time.Time, error) {

	// Convert periodDate to string and parse to time.Time
	dateStr := strconv.Itoa(periodDate)
	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %w", err)
	}

	// Convert periodStartTime to string and parse hours and minutes
	startTimeStr := fmt.Sprintf("%04d", periodTime) // Ensure it's 4 digits
	hours, _ := strconv.Atoi(startTimeStr[:2])
	minutes, _ := strconv.Atoi(startTimeStr[2:])
	// Load German timezone (CET/CEST)
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return time.Time{}, fmt.Errorf("error loading timezone: %w", err)
	}

	// Combine date with start time
	combinedTime := time.Date(date.Year(), date.Month(), date.Day(), hours, minutes, 0, 0, location)
	return combinedTime, nil
}

func (database *Database) FetchLesson(genUser gen.User, untis_pwd string, classId int, startDate time.Time, endDate time.Time, ctx context.Context) error {
	var user dbModels.User
	user.FromGen(genUser)
	err := database.fetchUser(&user, ctx)
	if err != nil {
		return err
	}
	periods, err := dataCollectors.DataCollectors.UntisClient.GetLessonsByStudent(user, untis_pwd, startDate, endDate, classId)
	if err != nil {
		return err
	}
	lessons := make([]dbModels.Lesson, len(periods))
	for i, period := range periods {
		var subjectIds []string
		for _, subject := range period.Subjects {
			subjectIds = append(subjectIds, fmt.Sprintf("%d", subject.ID))
		}
		var classIds []string
		for _, class := range period.Classes {
			classIds = append(classIds, fmt.Sprintf("%d", class.ID))
		}
		var teacherIds []string
		for _, teacher := range period.Teachers {
			teacherIds = append(teacherIds, fmt.Sprintf("%d", teacher.ID))
		}
		var roomIds []string
		for _, room := range period.Rooms {
			roomIds = append(roomIds, fmt.Sprintf("%d", room.ID))
		}
		startTime, err := MergeDateAndTime(period.Date, period.StartTime)
		if err != nil {
			return err
		}
		endTime, err := MergeDateAndTime(period.Date, period.EndTime)
		if err != nil {
			return err
		}
		substitutionText := period.SubstitutionText
		chairUp := false
		re := regexp.MustCompile(`Bitte aufstuhlen!*`) // TODO: To config
		if re.MatchString(substitutionText) {
			substitutionText = re.ReplaceAllString(substitutionText, "")
			chairUp = true
		}
		lessons[i] = dbModels.Lesson{
			Id:                    period.Id,
			Subjects:              subjectIds,
			Classes:               classIds,
			Teachers:              teacherIds,
			Rooms:                 roomIds,
			StartTime:             startTime,
			EndTime:               endTime,
			Cancelled:             (period.Code == "cancelled"),
			Irregular:             (period.Code == "irregular"),
			LessonType:            gen.LessonLessonType(period.LessonType),
			AdditionalInformation: period.Info,
			SubstitutionText:      substitutionText,
			LessonText:            period.LessonText,
			BookingText:           period.BookingText,
			ChairUp:               chairUp,
		}
	}
	lessonQuery := database.DB.NewInsert()
	lessonQuery.Model(&lessons)
	lessonQuery.On("CONFLICT (id) DO UPDATE")
	_, err = lessonQuery.Exec(ctx)
	return err
}
func placeholderArray(arr []string) string {
	placeholders := make([]string, len(arr))
	for i := range placeholders {
		placeholders[i] = "'?'"
	}
	// Join the placeholders and format the output
	return "array[" + strings.Join(placeholders, ", ") + "]"
}
func dataArray(arr []string) []interface{} {
	data := make([]interface{}, len(arr))
	for i, val := range arr {
		data[i] = bun.Safe(val)
	}
	// Join the placeholders and format the output
	return data
}

func parseInterfaceToStringArray(value interface{}) ([]string, error) {
	// Check if the provided value is a slice
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value is not a slice")
	}

	// Create a new []string slice with the same length as the input slice
	strArray := make([]string, v.Len())

	// Iterate over the slice and convert each element to a string
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()

		// Convert each element to string based on its type
		switch val := elem.(type) {
		case int:
			strArray[i] = strconv.Itoa(val)
		case float64:
			strArray[i] = strconv.FormatFloat(val, 'f', -1, 64)
		case bool:
			strArray[i] = strconv.FormatBool(val)
		case string:
			strArray[i] = val
		default:
			// Fallback to default string conversion for unsupported types
			strArray[i] = fmt.Sprintf("%v", elem)
		}
	}

	return strArray, nil
}

func (database *Database) GetLesson(filter dbModels.LessonFilter, ctx context.Context) ([]gen.Lesson, error) {
	if filter.User.Id == 0 {
		return nil, errors.New("user ID is required to get lessons")
	}
	var choice dbModels.Choice
	// get Choice
	if filter.Choice.Id == 0 {
		if filter.Choice.Choice == "" {
			query := database.DB.NewSelect()
			query.Model(&filter.User)
			query.WherePK()
			query.Relation("DefaultChoice")
			err := query.Scan(ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, dbModels.ErrUserNotFound
				}
				return nil, err
			}
			choice = *filter.User.DefaultChoice
		} else {
			choice = filter.Choice
			userQuery := database.DB.NewSelect()
			userQuery.Model(&filter.User)
			userQuery.WherePK()
			err := userQuery.Scan(ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, dbModels.ErrUserNotFound
				}
				return nil, err
			}
			//TODO: check class
		}

	} else if filter.Choice.Id != 0 {
		userQuery := database.DB.NewSelect()
		userQuery.Model(&filter.User)
		userQuery.WherePK()
		err := userQuery.Scan(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, dbModels.ErrUserNotFound
			}
			return nil, err
		}

		query := database.DB.NewSelect()
		query.Model(&filter.Choice)
		query.WherePK()
		if filter.User.Role != string(gen.UserRoleAdmin) {
			query.Where("\"choice\".\"userId\" = ?", filter.User.Id) //TODO: disable if Admin
		}
		err = query.Scan(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, dbModels.ErrChoiceNotFound
			}
			return nil, err
		}
		choice = filter.Choice
	}

	lessonQuery := database.DB.NewSelect()
	lessons := make([]dbModels.Lesson, 0)
	lessonQuery.Model(&lessons)

	var result map[string]interface{}
	parsingError := json.Unmarshal([]byte(choice.Choice), &result)

	if parsingError != nil || choice.Choice == "" || len(result) == 0 {
		lessonQuery.Where("\"lesson\".\"classes\" @> ?", pgdialect.Array(filter.User.Classes))
	} else {
		for key, value := range result {
			if classID, err := strconv.Atoi(key); err == nil {
				subjects, err := parseInterfaceToStringArray(value)
				if err != nil {
					return nil, err
				}
				if len(subjects) == 0 {
					lessonQuery.WhereOr("(\"lesson\".\"classes\" \\?| ARRAY[?])", key)
				} else if classID > 0 {
					lessonQuery.WhereOr("(\"lesson\".\"classes\" \\?| ARRAY[?] AND \"lesson\".\"subjects\" \\?| "+placeholderArray(subjects)+")", append([]interface{}{key}, dataArray(subjects)...)...)
				} else { //TODO: If a Class ID is present as a negative as well as a positive value only the positive should be used.
					lessonQuery.WhereOr("(\"lesson\".\"classes\" \\?| ARRAY[?] AND NOT \"lesson\".\"subjects\" \\?| "+placeholderArray(subjects)+")", append([]interface{}{key}, dataArray(subjects)...)...)

				}
			}
		}
	}
	if !filter.StartDate.IsZero() && !filter.EndDate.IsZero() {
		lessonQuery.Where("start_time >= ? AND end_time <= ?", filter.StartDate, filter.EndDate)
	}
	err := lessonQuery.Scan(ctx)
	genLesson := make([]gen.Lesson, len(lessons))
	for i, c := range lessons {
		genLesson[i] = c.ToGen()
	}
	return genLesson, err
}

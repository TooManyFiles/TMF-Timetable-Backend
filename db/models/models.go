package dbModels

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/uptrace/bun"
)

var ErrUsernameNotMachRequirements = errors.New("Username dose not match the requirements")
var ErrPasswordNotMachRequirements = errors.New("crypto: Password dose not match the requirements")
var ErrUserNotFound = errors.New("db: User not found")
var ErrInvalidPassword = errors.New("db: The Password is wrong")
var ErrInvalidToken = errors.New("db: The Token is invalid")
var ErrChoiceNotFound = errors.New("db: Choice not found")

func getPointerIfNotEmpty[T any](v T) *T {
	val := reflect.ValueOf(v)

	// Check if the value is zero (e.g., nil or empty)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		if val.Len() > 0 {
			return &v
		}
	case reflect.String:
		if strVal, ok := any(v).(string); ok && strVal != "" {
			return &v
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// For numeric types, return a pointer regardless of the value
		return &v
	default:
		if !val.IsZero() {
			return &v
		}
	}
	return nil
}

// Class model
type Class struct {
	bun.BaseModel          `bun:"table:classes,alias:c"`
	Id                     int    `bun:"id,pk,autoincrement,notnull"`
	Name                   string `bun:"name"`
	MainTeacherId          int    `bun:"mainTeacherId"`
	SecondaryTeacherId     int    `bun:"secondaryTeacherId"`
	MainClassLeaderId      int    `bun:"mainClassleader"`
	SecondaryClassLeaderId int    `bun:"secondaryClassleader"`

	MainTeacher          *Teacher `bun:"rel:belongs-to,join:mainTeacherId=id"`
	SecondaryTeacher     *Teacher `bun:"rel:belongs-to,join:secondaryTeacherId=id"`
	MainClassleader      *User    `bun:"rel:belongs-to,join:mainClassleader=id"`
	SecondaryClassleader *User    `bun:"rel:belongs-to,join:secondaryClassleader=id"`
}

func (class *Class) ToGen() gen.Class {
	return gen.Class{
		Id:                     getPointerIfNotEmpty(class.Id),
		Name:                   getPointerIfNotEmpty(class.Name),
		MainTeacherId:          getPointerIfNotEmpty(class.MainTeacherId),
		SecondaryTeacherId:     getPointerIfNotEmpty(class.SecondaryTeacherId),
		MainClassLeaderId:      getPointerIfNotEmpty(class.MainClassLeaderId),
		SecondaryClassLeaderId: getPointerIfNotEmpty(class.SecondaryClassLeaderId),
	}
}
func (class *Class) FromGen(genClass gen.Class) Class {
	if class == nil {
		class = &Class{}
	}
	if genClass.Id != nil {
		class.Id = int(*genClass.Id)
	}
	if *genClass.Name != "" {
		class.Name = *genClass.Name
	}
	if genClass.MainTeacherId != nil {
		class.MainTeacherId = *genClass.MainTeacherId
	}
	if genClass.SecondaryTeacherId != nil {
		class.SecondaryTeacherId = *genClass.SecondaryTeacherId
	}
	if genClass.MainClassLeaderId != nil {
		class.MainClassLeaderId = *genClass.MainClassLeaderId
	}
	if genClass.MainClassLeaderId != nil {
		class.SecondaryClassLeaderId = *genClass.MainClassLeaderId
	}
	return *class
}

// User model
type User struct {
	bun.BaseModel   `bun:"table:user"`
	Id              int      `bun:"id,pk,autoincrement,notnull"`
	Name            string   `bun:"name,unique"`
	Role            string   `bun:"role"`
	DefaultChoiceId int      `bun:"defaultChoice"`
	PwdHash         string   `bun:"pwdHash"`
	Classes         []string `pg:"classes,array"`
	Email           string   `pg:"email"`
	UntisPWD        string   `pg:"untispwd"`
	UntisName       string   `pg:"untisname"`
	UntisId         int      `pg:"untisid"`
	DefaultChoice   *Choice  `bun:"rel:belongs-to,join:defaultChoice=id"`
	Class           *Class   `bun:"rel:belongs-to,join:classes=id"`
}

func (user *User) FromGen(genUser gen.User) User {
	strClasses := make([]string, len(*genUser.Classes))
	for i, s := range *genUser.Classes {
		num := strconv.Itoa(s)
		strClasses[i] = num
	}
	if user == nil {
		user = &User{}
	}
	if genUser.Id != nil {
		user.Id = int(*genUser.Id)
	}
	if genUser.Name != "" {
		user.Name = genUser.Name
	}
	if genUser.Role != nil {
		user.Role = string(*genUser.Role)
	}
	if genUser.DefaultChoice != nil && genUser.DefaultChoice.Id != nil {
		user.DefaultChoiceId = *genUser.DefaultChoice.Id
	}
	if genUser.Classes != nil {

		user.Classes = strClasses
	}
	if genUser.Email != nil {
		user.Email = *genUser.Email
	}
	return *user
}
func (user *User) ToGen() gen.User {
	role := gen.UserRole(user.Role)

	intClasses := make([]int, len(user.Classes))
	for i, s := range user.Classes {
		num, err := strconv.Atoi(s)
		if err == nil {
			intClasses[i] = num
		}
	}

	if user.DefaultChoice != nil {
		choice := user.DefaultChoice.ToGen()
		return gen.User{
			Id:            getPointerIfNotEmpty(user.Id),
			Name:          user.Name,
			Role:          getPointerIfNotEmpty(role),
			DefaultChoice: getPointerIfNotEmpty(choice),
			Classes:       getPointerIfNotEmpty(intClasses),
			Email:         getPointerIfNotEmpty(user.Email),
		}
	}
	return gen.User{
		Id:            getPointerIfNotEmpty(user.Id),
		Name:          user.Name,
		Role:          getPointerIfNotEmpty(role),
		DefaultChoice: getPointerIfNotEmpty(gen.Choice{Id: &user.DefaultChoiceId}),
		Classes:       getPointerIfNotEmpty(intClasses),
		Email:         getPointerIfNotEmpty(user.Email),
	}
}

// Teacher model
type Teacher struct {
	bun.BaseModel `bun:"table:teacher"`
	Id            int `bun:"id,pk,autoincrement,notnull"`
	UserId        int `bun:"userId"`
	ShortName     string
	Name          string
	FirstName     string
	Pronoun       string
	Title         string
	User          *User `bun:"rel:belongs-to,join:userId=id"`
}

func (teacher *Teacher) ToGen() gen.Teacher {
	return gen.Teacher{
		Id:        getPointerIfNotEmpty(teacher.Id),
		UserId:    getPointerIfNotEmpty(teacher.UserId),
		Name:      getPointerIfNotEmpty(teacher.Name),
		FirstName: getPointerIfNotEmpty(teacher.FirstName),
		Pronoun:   getPointerIfNotEmpty(teacher.Pronoun),
		Title:     getPointerIfNotEmpty(teacher.Title),
		ShortName: getPointerIfNotEmpty(teacher.ShortName),
	}
}
func (teacher *Teacher) FromGen(genTeacher gen.Teacher) Teacher {
	if teacher == nil {
		teacher = &Teacher{}
	}
	if genTeacher.Id != nil {
		teacher.Id = int(*genTeacher.Id)
	}
	if genTeacher.UserId != nil {
		teacher.UserId = *genTeacher.UserId
	}
	if *genTeacher.Name != "" {
		teacher.Name = *genTeacher.Name
	}
	if *genTeacher.FirstName != "" {
		teacher.FirstName = *genTeacher.FirstName
	}
	if *genTeacher.Pronoun != "" {
		teacher.Pronoun = *genTeacher.Pronoun
	}
	if *genTeacher.Title != "" {
		teacher.Title = *genTeacher.Title
	}
	return *teacher
}

// Lesson model
type Lesson struct {
	bun.BaseModel         `bun:"table:lesson"`
	Id                    int       `bun:"id,pk,autoincrement,notnull"`
	Subjects              []string  `pg:",array"`
	Classes               []string  `pg:",array"`
	Teachers              []string  `pg:",array"`
	Rooms                 []string  `pg:",array"`
	OriginalSubjects      []string  `pg:",array"`
	OriginalClasses       []string  `pg:",array"`
	OriginalTeachers      []string  `pg:",array"`
	OriginalRooms         []string  `pg:",array"`
	StartTime             time.Time // Date-time format in Go can be parsed as time.Time
	EndTime               time.Time
	LastUpdate            time.Time
	Cancelled             bool   `json:"cancelled,omitempty"`
	Homework              string `json:"homework,omitempty"`
	Irregular             bool   `json:"irregular,omitempty"`
	ChairUp               bool
	AdditionalInformation string
	SubstitutionText      string
	LessonText            string
	BookingText           string
	// LessonType //„ls“ (lesson) | „oh“ (office hour) | „sb“ (standby) | „bs“ (break supervision) | „ex“(examination)  omitted if lesson
	LessonType gen.LessonLessonType `json:"lessonType"`
}

var _ bun.BeforeAppendModelHook = (*Lesson)(nil)

func (l *Lesson) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		l.LastUpdate = time.Now()
	case *bun.UpdateQuery:
		l.LastUpdate = time.Now()
	}
	return nil
}

func (lesson *Lesson) ToGen() gen.Lesson {
	intSubjects := make([]int, len(lesson.Subjects))
	for i, s := range lesson.Subjects {
		num, err := strconv.Atoi(s)
		if err == nil {
			intSubjects[i] = num
		}
	}
	intClasses := make([]int, len(lesson.Classes))
	for i, s := range lesson.Classes {
		num, err := strconv.Atoi(s)
		if err == nil {
			intClasses[i] = num
		}
	}
	intTeachers := make([]int, len(lesson.Teachers))
	for i, s := range lesson.Teachers {
		num, err := strconv.Atoi(s)
		if err == nil {
			intTeachers[i] = num
		}
	}
	intRooms := make([]int, len(lesson.Rooms))
	for i, s := range lesson.Rooms {
		num, err := strconv.Atoi(s)
		if err == nil {
			intRooms[i] = num
		}
	}
	intOriginalSubjects := make([]int, len(lesson.OriginalSubjects))
	for i, s := range lesson.OriginalSubjects {
		num, err := strconv.Atoi(s)
		if err == nil {
			intOriginalSubjects[i] = num
		}
	}
	intOriginalClasses := make([]int, len(lesson.OriginalClasses))
	for i, s := range lesson.OriginalClasses {
		num, err := strconv.Atoi(s)
		if err == nil {
			intOriginalClasses[i] = num
		}
	}
	intOriginalTeachers := make([]int, len(lesson.OriginalTeachers))
	for i, s := range lesson.OriginalTeachers {
		num, err := strconv.Atoi(s)
		if err == nil {
			intOriginalTeachers[i] = num
		}
	}
	intOriginalRooms := make([]int, len(lesson.OriginalRooms))
	for i, s := range lesson.OriginalRooms {
		num, err := strconv.Atoi(s)
		if err == nil {
			intOriginalRooms[i] = num
		}
	}

	return gen.Lesson{
		Id:                    getPointerIfNotEmpty(lesson.Id),
		Subjects:              getPointerIfNotEmpty(intSubjects),
		Classes:               getPointerIfNotEmpty(intClasses),
		Teachers:              getPointerIfNotEmpty(intTeachers),
		Rooms:                 getPointerIfNotEmpty(intRooms),
		OrigSubjects:          getPointerIfNotEmpty(intOriginalSubjects),
		OrigClasses:           getPointerIfNotEmpty(intOriginalClasses),
		OrigTeachers:          getPointerIfNotEmpty(intOriginalTeachers),
		OrigRooms:             getPointerIfNotEmpty(intOriginalRooms),
		StartTime:             lesson.StartTime,
		EndTime:               lesson.EndTime,
		LastUpdate:            getPointerIfNotEmpty(lesson.LastUpdate),
		Cancelled:             getPointerIfNotEmpty(lesson.Cancelled),
		Irregular:             getPointerIfNotEmpty(lesson.Irregular),
		LessonType:            lesson.LessonType,
		AdditionalInformation: getPointerIfNotEmpty(lesson.AdditionalInformation),
		BookingText:           getPointerIfNotEmpty(lesson.BookingText),
		LessonText:            getPointerIfNotEmpty(lesson.LessonText),
		SubstitutionText:      getPointerIfNotEmpty(lesson.SubstitutionText),
		Homework:              getPointerIfNotEmpty(lesson.Homework),
		ChairUp:               getPointerIfNotEmpty(lesson.ChairUp),
	}
}
func (lesson *Lesson) FromGen(genLesson gen.Lesson) Lesson {
	// Convert Subjects
	strSubjects := make([]string, len(*genLesson.Subjects))
	for i, s := range *genLesson.Subjects {
		num := strconv.Itoa(s)
		strSubjects[i] = num

	}

	// Convert Classes
	strClasses := make([]string, len(*genLesson.Classes))
	for i, s := range *genLesson.Classes {
		num := strconv.Itoa(s)
		strClasses[i] = num

	}

	// Convert Teachers
	strTeachers := make([]string, len(*genLesson.Teachers))
	for i, s := range *genLesson.Teachers {
		num := strconv.Itoa(s)
		strTeachers[i] = num

	}

	// Convert Rooms
	strRooms := make([]string, len(*genLesson.Rooms))
	for i, s := range *genLesson.Rooms {
		num := strconv.Itoa(s)
		strRooms[i] = num

	}
	// Convert Subjects
	strOriginalSubjects := make([]string, len(*genLesson.OrigSubjects))
	for i, s := range *genLesson.OrigSubjects {
		num := strconv.Itoa(s)
		strOriginalSubjects[i] = num

	}

	// Convert Classes
	strOriginalClasses := make([]string, len(*genLesson.OrigClasses))
	for i, s := range *genLesson.OrigClasses {
		num := strconv.Itoa(s)
		strOriginalClasses[i] = num

	}

	// Convert Teachers
	strOriginalTeachers := make([]string, len(*genLesson.OrigTeachers))
	for i, s := range *genLesson.OrigTeachers {
		num := strconv.Itoa(s)
		strOriginalTeachers[i] = num

	}

	// Convert Rooms
	strOriginalRooms := make([]string, len(*genLesson.OrigRooms))
	for i, s := range *genLesson.OrigRooms {
		num := strconv.Itoa(s)
		strOriginalRooms[i] = num

	}
	if lesson == nil {
		lesson = &Lesson{}
	}
	if genLesson.Id != nil {
		lesson.Id = *genLesson.Id
	}
	if genLesson.Subjects != nil {
		lesson.Subjects = strSubjects
	}
	if genLesson.Classes != nil {
		lesson.Classes = strClasses
	}
	if genLesson.Teachers != nil {
		lesson.Teachers = strTeachers
	}
	if genLesson.Rooms != nil {
		lesson.Rooms = strRooms
	}
	if genLesson.OrigSubjects != nil {
		lesson.OriginalSubjects = strOriginalSubjects
	}
	if genLesson.OrigClasses != nil {
		lesson.OriginalClasses = strOriginalClasses
	}
	if genLesson.OrigTeachers != nil {
		lesson.OriginalTeachers = strOriginalTeachers
	}
	if genLesson.OrigRooms != nil {
		lesson.OriginalRooms = strOriginalRooms
	}
	if genLesson.StartTime.IsZero() {
		lesson.StartTime = genLesson.StartTime
	}
	if genLesson.EndTime.IsZero() {
		lesson.EndTime = genLesson.EndTime
	}
	if genLesson.LastUpdate != nil {
		lesson.LastUpdate = *genLesson.LastUpdate
	}
	if genLesson.Cancelled != nil {
		lesson.Cancelled = *genLesson.Cancelled
	}
	if genLesson.Irregular != nil {
		lesson.Irregular = *genLesson.Irregular
	}
	if genLesson.LessonType != "" {
		lesson.LessonType = genLesson.LessonType
	}
	if genLesson.AdditionalInformation != nil {
		lesson.AdditionalInformation = *genLesson.AdditionalInformation
	}
	if genLesson.BookingText != nil {
		lesson.BookingText = *genLesson.BookingText
	}
	if genLesson.LessonText != nil {
		lesson.LessonText = *genLesson.LessonText
	}
	if genLesson.SubstitutionText != nil {
		lesson.SubstitutionText = *genLesson.SubstitutionText
	}
	if genLesson.Homework != nil {
		lesson.Homework = *genLesson.Homework
	}
	if genLesson.ChairUp != nil {
		lesson.ChairUp = *genLesson.ChairUp
	}
	return *lesson
}

// Room model
type Room struct {
	bun.BaseModel         `bun:"table:room"`
	Id                    int `bun:"id,pk,autoincrement,notnull"`
	Name                  string
	AdditionalInformation string
}

func (room *Room) ToGen() gen.Room {
	return gen.Room{
		Id:                    getPointerIfNotEmpty(room.Id),
		Name:                  getPointerIfNotEmpty(room.Name),
		AdditionalInformation: getPointerIfNotEmpty(room.AdditionalInformation),
	}
}
func (room *Room) FromGen(genRoom gen.Room) Room {
	if room == nil {
		room = &Room{}
	}
	if genRoom.Id != nil {
		room.Id = int(*genRoom.Id)
	}
	if *genRoom.Name != "" {
		room.Name = *genRoom.Name
	}
	if *genRoom.AdditionalInformation != "" {
		room.AdditionalInformation = *genRoom.AdditionalInformation
	}

	return *room
}

// Subject model
type Subject struct {
	bun.BaseModel `bun:"table:subject"`
	Id            int `bun:"id,pk,autoincrement,notnull"`
	Name          string
	ShortName     string
}

func (subject *Subject) ToGen() gen.Subject {
	return gen.Subject{
		Id:        getPointerIfNotEmpty(subject.Id),
		Name:      getPointerIfNotEmpty(subject.Name),
		ShortName: getPointerIfNotEmpty(subject.ShortName),
	}
}
func (subject *Subject) FromGen(genSubject gen.Subject) Subject {
	if subject == nil {
		subject = &Subject{}
	}
	if genSubject.Id != nil {
		subject.Id = int(*genSubject.Id)
	}
	if *genSubject.Name != "" {
		subject.Name = *genSubject.Name
	}
	if *genSubject.ShortName != "" {
		subject.ShortName = *genSubject.ShortName
	}
	return *subject
}

type Choice struct {
	bun.BaseModel `bun:"table:choice"`
	Id            int `bun:"id,pk,autoincrement,notnull"`
	UserId        int `bun:"userId"`
	Name          string
	Choice        string // Assuming this is a JSON field
	User          *User  `bun:"rel:belongs-to,join:userId=id"`
}

func (choice *Choice) ToGen() gen.Choice {
	var choiceMap map[string]interface{}

	// Unmarshal the JSON string into a map
	err := json.Unmarshal([]byte(choice.Choice), &choiceMap)
	if err != nil {
		return gen.Choice{
			Id:     getPointerIfNotEmpty(choice.Id),
			Name:   getPointerIfNotEmpty(choice.Name),
			UserId: getPointerIfNotEmpty(choice.UserId),
		}
	}
	return gen.Choice{
		Id:     getPointerIfNotEmpty(choice.Id),
		Name:   getPointerIfNotEmpty(choice.Name),
		UserId: getPointerIfNotEmpty(choice.UserId),
		Choice: getPointerIfNotEmpty(choiceMap),
	}
}
func (choice *Choice) FromGen(genChoice gen.Choice) Choice {
	if choice == nil {
		choice = &Choice{}
	}
	if genChoice.Id != nil {
		choice.Id = int(*genChoice.Id)
	}
	if genChoice.Name != nil && *genChoice.Name != "" {
		choice.Name = *genChoice.Name
	}
	if genChoice.UserId != nil {
		choice.UserId = *genChoice.UserId
	}
	jsonChoice, _ := json.Marshal(genChoice.Choice)
	if genChoice.Choice != nil {
		choice.Choice = string(jsonChoice)
	}
	return *choice
}

type LessonFilter struct {
	Choice    Choice
	User      User
	StartDate time.Time
	EndDate   time.Time
}
type Menu struct {
	bun.BaseModel `bun:"table:menu"`
	Date          time.Time `bun:"date,pk,unique,notnull,type:date" json:"date,omitempty"`
	Cookteam      string    `json:"cookteam,omitempty"`
	Dessert       string    `json:"dessert,omitempty"`
	Garnish       string    `json:"garnish,omitempty"`
	MainDish      string    `json:"mainDish,omitempty"`
	MainDishVeg   string    `json:"mainDishVeg,omitempty"`
	UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	NotAPIServed  bool      `bun:"notAPIServed,notnull,default:false"`
}

var _ bun.BeforeAppendModelHook = (*Menu)(nil)

func (u *Menu) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		u.UpdatedAt = time.Now()
	}
	return nil
}

package dbModels

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/uptrace/bun"
)

var ErrPasswordNotMachRequirements = errors.New("crypto: Password dose not match the requirements")
var ErrUserNotFound = errors.New("db: User not found")
var ErrInvalidPassword = errors.New("db: The Password is wrong")
var ErrInvalidToken = errors.New("db: The Token is invalid")

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
		Id:                     &class.Id,
		Name:                   &class.Name,
		MainTeacherId:          &class.MainTeacherId,
		SecondaryTeacherId:     &class.SecondaryTeacherId,
		MainClassLeaderId:      &class.MainClassLeaderId,
		SecondaryClassLeaderId: &class.SecondaryClassLeaderId,
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
	Id              int     `bun:"id,pk,autoincrement,notnull"`
	Name            string  `bun:"name,unique"`
	Role            string  `bun:"role"`
	DefaultChoiceId int     `bun:"defaultChoice"`
	PwdHash         string  `bun:"pwdHash"`
	Classes         []int   `pg:"classes,array"`
	Email           string  `pg:"email"`
	DefaultChoice   *Choice `bun:"rel:belongs-to,join:defaultChoice=id"`
	Class           *Class  `bun:"rel:belongs-to,join:classes=id"`
}

func (user *User) FromGen(genUser gen.User) User {
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
		user.Classes = *genUser.Classes
	}
	if genUser.Email != nil {
		user.Email = *genUser.Email
	}
	return *user
}
func (user *User) ToGen() gen.User {
	role := gen.UserRole(user.Role)
	if user.DefaultChoice != nil {
		choice := user.DefaultChoice.ToGen()
		return gen.User{
			Id:            &user.Id,
			Name:          user.Name,
			Role:          &role,
			DefaultChoice: &choice,
			Classes:       &user.Classes,
			Email:         &user.Email,
		}
	}
	return gen.User{
		Id:            &user.Id,
		Name:          user.Name,
		Role:          &role,
		DefaultChoice: &gen.Choice{Id: &user.DefaultChoiceId},
		Classes:       &user.Classes,
		Email:         &user.Email,
	}
}

// Teacher model
type Teacher struct {
	bun.BaseModel `bun:"table:teacher"`
	Id            int `bun:"id,pk,autoincrement,notnull"`
	UserId        int `bun:"userId"`
	Name          string
	FirstName     string
	Pronoun       string
	Title         string
	User          *User `bun:"rel:belongs-to,join:userId=id"`
}

func (teacher *Teacher) ToGen() gen.Teacher {
	return gen.Teacher{
		Id:        &teacher.Id,
		UserId:    &teacher.UserId,
		Name:      &teacher.Name,
		FirstName: &teacher.FirstName,
		Pronoun:   &teacher.Pronoun,
		Title:     &teacher.Title,
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
	bun.BaseModel `bun:"table:lesson"`
	Id            int    `bun:"id,pk,autoincrement,notnull"`
	Subjects      []int  `pg:",array"`
	Classes       []int  `pg:",array"`
	Teachers      []int  `pg:",array"`
	Rooms         []int  `pg:",array"`
	StartTime     string // Date-time format in Go can be parsed as time.Time
	EndTime       string
	LastUpdate    string
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
		Id:                    &room.Id,
		Name:                  &room.Name,
		AdditionalInformation: &room.AdditionalInformation,
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
		Id:        &subject.Id,
		Name:      &subject.Name,
		ShortName: &subject.ShortName,
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
			Id:     &choice.Id,
			Name:   &choice.Name,
			UserId: &choice.UserId,
		}
	}
	return gen.Choice{
		Id:     &choice.Id,
		Name:   &choice.Name,
		UserId: &choice.UserId,
		Choice: &choiceMap,
	}
}
func (choice *Choice) FromGen(genChoice gen.Choice) Choice {
	if choice == nil {
		choice = &Choice{}
	}
	if genChoice.Id != nil {
		choice.Id = int(*genChoice.Id)
	}
	if *genChoice.Name != "" {
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

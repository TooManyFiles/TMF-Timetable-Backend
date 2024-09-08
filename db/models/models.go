package dbModels

import (
	"context"
	"errors"
	"time"

	"github.com/TooManyFiles/TMF-Timetable-Backend/api/gen"
	"github.com/uptrace/bun"
)

var ErrPasswordNotMachRequirements = errors.New("crypto: Password dose not match the requirements.")
var ErrUserNotFound = errors.New("db: User not found.")
var ErrInvalidPassword = errors.New("db: The Password is wrong.")

// Class model
type Class struct {
	bun.BaseModel        `bun:"table:classes,alias:c"`
	Id                   int    `bun:"id,pk,autoincrement,notnull"`
	Name                 string `bun:"name"`
	MainTeacherId        int    `bun:"mainTeacherId"`
	SecondaryTeacherId   int    `bun:"secondaryTeacherId"`
	MainClassleader      int    `bun:"mainClassleader"`
	SecondaryClassleader int    `bun:"secondaryClassleader"`

	Teacher *Teacher `bun:"rel:belongs-to,join:mainTeacherId=id,join:secondaryTeacherId=id"`
	User    *User    `bun:"rel:belongs-to,join:mainClassleader=id,join:secondaryClassleader=id"`
}

// User model
type User struct {
	bun.BaseModel `bun:"table:user"`
	Id            int     `bun:"id,pk,autoincrement,notnull"`
	Name          string  `bun:"name,unique"`
	Role          string  `bun:"role"`
	DefaultChoice int     `bun:"defaultChoice"`
	PwdHash       string  `bun:"pwdHash"`
	Classes       []int   `pg:"classes,array"`
	Email         string  `pg:"email"`
	Choice        *Choice `bun:"rel:belongs-to,join:defaultChoice=id"`
	Class         *Class  `bun:"rel:belongs-to,join:classes=id"`
}

func (user *User) FromGenUser(genUser gen.User) User {
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
		user.DefaultChoice = *genUser.DefaultChoice.Id
	}
	if genUser.Classes != nil {
		user.Classes = *genUser.Classes
	}
	if genUser.Email != nil {
		user.Email = *genUser.Email
	}
	return *user
}
func (user *User) ToGenUser() gen.User {
	role := gen.UserRole(user.Role)
	return gen.User{
		Id:            &user.Id,
		Name:          user.Name,
		Role:          &role,
		DefaultChoice: &gen.Choice{Id: &user.DefaultChoice},
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

// Subject model
type Subject struct {
	bun.BaseModel `bun:"table:subject"`
	Id            int `bun:"id,pk,autoincrement,notnull"`
	Name          string
}
type Choice struct {
	bun.BaseModel `bun:"table:choice"`
	Id            int `bun:"id,pk,autoincrement,notnull"`
	UserId        int `bun:"userId"`
	Name          string
	Choice        string // Assuming this is a JSON field
	User          *User  `bun:"rel:belongs-to,join:userId=id"`
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

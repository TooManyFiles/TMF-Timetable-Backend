package dbModels

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Class model
type Class struct {
	bun.BaseModel        `bun:"table:classes,alias:c"`
	Id                   int64  `bun:"id,pk,autoincrement,notnull"`
	Name                 string `bun:"name"`
	MainTeacherId        int64  `bun:"mainTeacherId"`
	SecondaryTeacherId   int64  `bun:"secondaryTeacherId"`
	MainClassleader      int64  `bun:"mainClassleader"`
	SecondaryClassleader int64  `bun:"secondaryClassleader"`

	Teacher *Teacher `bun:"rel:belongs-to,join:mainTeacherId=id,join:secondaryTeacherId=id"`
	User    *User    `bun:"rel:belongs-to,join:mainClassleader=id,join:secondaryClassleader=id"`
}

// User model
type User struct {
	bun.BaseModel `bun:"table:user"`
	Id            int64   `bun:"id,pk,autoincrement,notnull"`
	Name          string  `bun:"name"`
	Role          string  `bun:"role"`
	DefaultChoice string  `bun:"defaultChoice"`
	PwdHash       string  `bun:"pwdHash"`
	Classes       []int64 `pg:"classes,array"`
	Choice        *Choice `bun:"rel:belongs-to,join:defaultChoice=id"`
	Class         *Class  `bun:"rel:belongs-to,join:classes=id"`
}

// Teacher model
type Teacher struct {
	bun.BaseModel `bun:"table:teacher"`
	Id            int64 `bun:"id,pk,autoincrement,notnull"`
	UserId        int64 `bun:"userId"`
	Name          string
	FirstName     string
	Pronoun       string
	Title         string
	User          *User `bun:"rel:belongs-to,join:userId=id"`
}

// Lesson model
type Lesson struct {
	bun.BaseModel `bun:"table:lesson"`
	Id            int64   `bun:"id,pk,autoincrement,notnull"`
	Subjects      []int64 `pg:",array"`
	Classes       []int64 `pg:",array"`
	Teachers      []int64 `pg:",array"`
	Rooms         []int64 `pg:",array"`
	StartTime     string  // Date-time format in Go can be parsed as time.Time
	EndTime       string
	LastUpdate    string
}

// Room model
type Room struct {
	bun.BaseModel         `bun:"table:room"`
	Id                    int64 `bun:"id,pk,autoincrement,notnull"`
	Name                  string
	AdditionalInformation string
}

// Subject model
type Subject struct {
	bun.BaseModel `bun:"table:subject"`
	Id            int64 `bun:"id,pk,autoincrement,notnull"`
	Name          string
}
type Choice struct {
	bun.BaseModel `bun:"table:choice"`
	Id            int64 `bun:"id,pk,autoincrement,notnull"`
	UserId        int64 `bun:"userId"`
	Name          string
	Choice        string // Assuming this is a JSON field
	User          *User  `bun:"rel:belongs-to,join:userId=id"`
}

type Menu struct {
	bun.BaseModel `bun:"table:menu"`
	Id            int64     `bun:"id,pk,autoincrement,notnull"`
	Cookteam      string    `json:"cookteam,omitempty"`
	Date          time.Time `bun:"date,unique,notnull" json:"date,omitempty"`
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

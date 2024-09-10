// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package gen

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for LessonLessonType.
const (
	Bs LessonLessonType = "bs"
	Ex LessonLessonType = "ex"
	Ls LessonLessonType = "ls"
	Oh LessonLessonType = "oh"
	Sb LessonLessonType = "sb"
)

// Defines values for UserRole.
const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleStudent UserRole = "student"
	UserRoleTeacher UserRole = "teacher"
)

// Defines values for PutViewJSONBodyProvider.
const (
	PutViewJSONBodyProviderCafeteria PutViewJSONBodyProvider = "cafeteria"
	PutViewJSONBodyProviderUntis     PutViewJSONBodyProvider = "untis"
)

// Defines values for PutViewUserUserIdJSONBodyProvider.
const (
	PutViewUserUserIdJSONBodyProviderCafeteria PutViewUserUserIdJSONBodyProvider = "cafeteria"
	PutViewUserUserIdJSONBodyProviderUntis     PutViewUserUserIdJSONBodyProvider = "untis"
)

// Choice Choice of subjects for the classes. {class:[subjects]}
// - If a class has a empty array as a choice all subjects should be shown.
// - If the Class ID is negative it the the choice is a blacklist.
// - If a Class ID is present as a negative as well as a positive value only the positive should be used.
type Choice struct {
	Choice *map[string]interface{} `json:"Choice,omitempty"`
	Id     *int                    `json:"id,omitempty"`
	Name   *string                 `json:"name,omitempty"`
	UserId *int                    `json:"userId,omitempty"`
}

// Class defines model for Class.
type Class struct {
	Id                     *int    `json:"id,omitempty"`
	MainClassLeaderID      *int    `json:"mainClassLeaderID,omitempty"`
	MainTeacherId          *int    `json:"mainTeacherId,omitempty"`
	Name                   *string `json:"name,omitempty"`
	SecondaryClassLeaderID *int    `json:"secondaryClassLeaderID,omitempty"`
	SecondaryTeacherId     *int    `json:"secondaryTeacherId,omitempty"`
}

// Lesson defines model for Lesson.
type Lesson struct {
	AdditionalInformation *string    `json:"additionalInformation,omitempty"`
	Cancelled             *bool      `json:"cancelled,omitempty"`
	Classes               *[]int     `json:"classes,omitempty"`
	EndTime               time.Time  `json:"endTime"`
	Homework              *string    `json:"homework,omitempty"`
	Id                    *int       `json:"id,omitempty"`
	Irregular             *bool      `json:"irregular,omitempty"`
	LastUpdate            *time.Time `json:"lastUpdate,omitempty"`

	// LessonType //„ls“ (lesson) | „oh“ (office hour) | „sb“ (standby) | „bs“ (break supervision) | „ex“(examination)  omitted if lesson
	LessonType LessonLessonType `json:"lessonType"`
	Rooms      *[]int           `json:"rooms,omitempty"`
	StartTime  time.Time        `json:"startTime"`
	Subjects   *[]int           `json:"subjects,omitempty"`
	Teachers   *[]int           `json:"teachers,omitempty"`
}

// LessonLessonType //„ls“ (lesson) | „oh“ (office hour) | „sb“ (standby) | „bs“ (break supervision) | „ex“(examination)  omitted if lesson
type LessonLessonType string

// Menu defines model for Menu.
type Menu struct {
	Cookteam    *string            `json:"cookteam,omitempty"`
	Date        openapi_types.Date `json:"date"`
	Dessert     *string            `json:"dessert,omitempty"`
	Garnish     *string            `json:"garnish,omitempty"`
	MainDish    *string            `json:"mainDish,omitempty"`
	MainDishVeg *string            `json:"mainDishVeg,omitempty"`
}

// Room defines model for Room.
type Room struct {
	AdditionalInformation *string `json:"additionalInformation,omitempty"`
	Id                    *int    `json:"id,omitempty"`
	Name                  *string `json:"name,omitempty"`
}

// Subject defines model for Subject.
type Subject struct {
	Id        *int    `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	ShortName *string `json:"shortName,omitempty"`
}

// Teacher defines model for Teacher.
type Teacher struct {
	FirstName *string `json:"firstName,omitempty"`
	Id        *int    `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	Pronoun   *string `json:"pronoun,omitempty"`
	Title     *string `json:"title,omitempty"`
	UserId    *int    `json:"userId,omitempty"`
}

// User defines model for User.
type User struct {
	Classes *[]int `json:"classes,omitempty"`

	// DefaultChoice Choice of subjects for the classes. {class:[subjects]}
	// - If a class has a empty array as a choice all subjects should be shown.
	// - If the Class ID is negative it the the choice is a blacklist.
	// - If a Class ID is present as a negative as well as a positive value only the positive should be used.
	DefaultChoice *Choice   `json:"defaultChoice,omitempty"`
	Email         *string   `json:"email,omitempty"`
	Id            *int      `json:"id,omitempty"`
	Name          string    `json:"name"`
	Role          *UserRole `json:"role,omitempty"`
}

// UserRole defines model for User.Role.
type UserRole string

// GetCafeteriaParams defines parameters for GetCafeteria.
type GetCafeteriaParams struct {
	Date     *openapi_types.Date `form:"date,omitempty" json:"date,omitempty"`
	Duration *int                `form:"duration,omitempty" json:"duration,omitempty"`
}

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody struct {
	Email *string `json:"email,omitempty"`

	// Password yourpassword hashed with SHA256
	Password *string `json:"password,omitempty"`
}

// PutViewJSONBody defines parameters for PutView.
type PutViewJSONBody struct {
	Provider []PutViewJSONBodyProvider `json:"provider"`
	Untis    *struct {
		// Choice Choice of subjects for the classes. {class:[subjects]}
		// - If a class has a empty array as a choice all subjects should be shown.
		// - If the Class ID is negative it the the choice is a blacklist.
		// - If a Class ID is present as a negative as well as a positive value only the positive should be used.
		Choice *Choice `json:"Choice,omitempty"`
	} `json:"untis,omitempty"`
}

// PutViewParams defines parameters for PutView.
type PutViewParams struct {
	Date     *openapi_types.Date `form:"date,omitempty" json:"date,omitempty"`
	Duration *int                `form:"duration,omitempty" json:"duration,omitempty"`
}

// PutViewJSONBodyProvider defines parameters for PutView.
type PutViewJSONBodyProvider string

// PutViewUserUserIdJSONBody defines parameters for PutViewUserUserId.
type PutViewUserUserIdJSONBody struct {
	Provider []PutViewUserUserIdJSONBodyProvider `json:"provider"`
	Untis    *struct {
		// Choice Choice of subjects for the classes. {class:[subjects]}
		// - If a class has a empty array as a choice all subjects should be shown.
		// - If the Class ID is negative it the the choice is a blacklist.
		// - If a Class ID is present as a negative as well as a positive value only the positive should be used.
		Choice *Choice `json:"Choice,omitempty"`
	} `json:"untis,omitempty"`
}

// PutViewUserUserIdParams defines parameters for PutViewUserUserId.
type PutViewUserUserIdParams struct {
	Date     *openapi_types.Date `form:"date,omitempty" json:"date,omitempty"`
	Duration *int                `form:"duration,omitempty" json:"duration,omitempty"`
}

// PutViewUserUserIdJSONBodyProvider defines parameters for PutViewUserUserId.
type PutViewUserUserIdJSONBodyProvider string

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody = User

// PutUsersUserIdJSONRequestBody defines body for PutUsersUserId for application/json ContentType.
type PutUsersUserIdJSONRequestBody = User

// PostUsersUserIdChoicesChoiceIdJSONRequestBody defines body for PostUsersUserIdChoicesChoiceId for application/json ContentType.
type PostUsersUserIdChoicesChoiceIdJSONRequestBody = Choice

// PutViewJSONRequestBody defines body for PutView for application/json ContentType.
type PutViewJSONRequestBody PutViewJSONBody

// PutViewUserUserIdJSONRequestBody defines body for PutViewUserUserId for application/json ContentType.
type PutViewUserUserIdJSONRequestBody PutViewUserUserIdJSONBody

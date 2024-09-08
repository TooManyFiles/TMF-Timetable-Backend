package api

import "github.com/TooManyFiles/TMF-Timetable-Backend/db"

// optional code omitted

type Server struct {
	DB db.Database
}

func NewServer(DB db.Database) Server {
	return Server{
		DB: DB,
	}
}

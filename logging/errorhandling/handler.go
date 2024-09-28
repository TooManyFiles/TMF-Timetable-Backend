package errorhandling

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/TooManyFiles/TMF-Timetable-Backend/logging"
)

type ErrorLevel int

const (
	ErrorWARN ErrorLevel = iota
	ErrorWrongUsage
	ErrorMedium
	ErrorIntense
)

var Logger *logging.Logger

type ErrorSource struct {
	Name       string
	ParseError func(error) (CustomError, error)
}

func generateTraceId() string {
	length := 16                    // bytes
	bytes := make([]byte, length/2) // half the length because hex encoding doubles the size
	if _, err := rand.Read(bytes); err != nil {
		(&CustomError{CustomErrorPreset: FailedToGenerateTraceId, TraceId: generateTraceId()}).Log()
		return ""
	}
	return hex.EncodeToString(bytes)
}

type CustomError struct {
	TraceId string
	CustomErrorPreset
}

func NewCustomError(preset CustomErrorPreset) CustomError {
	return CustomError{CustomErrorPreset: preset, TraceId: generateTraceId()}
}

type CustomErrorPreset struct {
	Code        int
	UserMessage string
	DevMessage  string
	LogMessage  string
	Source      ErrorSource
	Level       ErrorLevel
	HttpCode    int
}

type HtmlError struct {
	UserMessage string `json:"userMessage"`
	DevMessage  string `json:"devMessage"`
	TraceId     string `json:"TraceId"`
}

// Implement the Error() method for CustomError to satisfy the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.LogMessage)
}

func (e *CustomError) HTML() (HtmlError, int) {
	return HtmlError{
		UserMessage: e.UserMessage,
		DevMessage:  e.DevMessage,
		TraceId:     e.TraceId,
	}, e.HttpCode
}
func (e *CustomError) Log() *CustomError {
	if e.Level == ErrorWARN {
		Logger.Warn(fmt.Sprintf("[%s]\tTrace: %s\t%s", e.Source.Name, e.TraceId, e.LogMessage))
	} else if e.Level == ErrorWrongUsage {
		Logger.Info(fmt.Sprintf("[%s]\tTrace: %s\t%s", e.Source.Name, e.TraceId, e.LogMessage))
	} else if e.Level == ErrorMedium {
		Logger.Error(fmt.Sprintf("[%s]\tTrace: %s\t%s", e.Source.Name, e.TraceId, e.LogMessage))
	} else if e.Level == ErrorIntense {
		Logger.Fail(fmt.Sprintf("[%s]\tTrace: %s\t%s", e.Source.Name, e.TraceId, e.LogMessage))
	} else {
		Logger.Error(fmt.Sprintf("[%s]\tTrace: %s\t%s", e.Source.Name, e.TraceId, e.LogMessage))
	}
	return e
}
func NewError(err error) CustomError {
	return CustomError{}
}

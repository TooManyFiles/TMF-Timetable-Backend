package logging

import (
	"fmt"
	"log"
	"strings"
)

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FAIL
	NONE
)

type LogLevel int

const (
	// ANSI escape codes for foreground colors
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	DarkRed = "\033[35m"

	// ANSI escape codes for background colors
	RedBG     = "\033[41m"
	GreenBG   = "\033[42m"
	YellowBG  = "\033[43m"
	BlueBG    = "\033[44m"
	DarkRedBG = "\033[45m"
)

// Logger structure
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// New logger constructor
func newLogger(logLogger *log.Logger, level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: logLogger,
	}
}

// Helper function to log messages with color and level
func (l *Logger) logWithLevel(color, level, textColor string, msg ...string) {
	if l.logger == nil {
		return
	}

	// Save the original prefix
	originalPrefix := l.logger.Prefix()
	if textColor == "" {
		textColor = Reset
	}
	// Set new prefix with color and level
	l.logger.SetPrefix(fmt.Sprintf("%s[%s]\t%s", color, level, textColor))

	// Log the message
	l.logger.Println(strings.Join(msg, " "))

	// Restore the original prefix
	l.logger.SetPrefix(originalPrefix)
}

// Debug level log with blue color
func (l *Logger) Debug(msg ...string) {
	if l.level <= DEBUG {
		l.logWithLevel(Blue, "DEBUG", "", msg...)
	}
}

// Info level log with green color
func (l *Logger) Info(msg ...string) {
	if l.level <= INFO {
		l.logWithLevel(Green, "INFO", "", msg...)
	}
}

// Warn level log with yellow color
func (l *Logger) Warn(msg ...string) {
	if l.level <= WARN {
		l.logWithLevel(Yellow, "WARN", "", msg...)
	}
}

// Error level log with red color
func (l *Logger) Error(msg ...string) {
	if l.level <= ERROR {
		l.logWithLevel(Red, "ERROR", "", msg...)
	}
}

// Error level log with red color
func (l *Logger) Fail(msg ...string) {
	if l.level <= FAIL {
		l.logWithLevel(Red+DarkRedBG, "FAIL", Red+DarkRedBG, msg...)
	}
}

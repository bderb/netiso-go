package logger

import (
	"log"
	"os"
)

const (
	ErrorLevel  = iota // Logs only errors, exits the program on error
	WarningLevel       // Logs warnings and errors
	InfoLevel          // Logs info, warnings, and errors
	DebugLevel         // Logs everything: debug, info, warning, error
)

type Logger struct {
	level int
}

func NewLogger(level int) *Logger {
	return &Logger{level: level}
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) GetLevel() int {
	return l.level
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level >= ErrorLevel {
		log.Printf("\r[ERROR] "+format, args...)
		os.Exit(1)
	}
}

func (l *Logger) Warning(format string, args ...interface{}) {
	if l.level >= ErrorLevel {
		log.Printf("\r[WARNING] "+format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level >= InfoLevel {
		log.Printf("\r[INFO] "+format, args...)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level >= DebugLevel {
		log.Printf("\r[DEBUG] "+format, args...)
	}
}
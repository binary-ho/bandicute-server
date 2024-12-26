package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Setup(level string) {
	// Set logging format
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
	})

	// Set output to stdout
	logrus.SetOutput(os.Stdout)

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Warn("Invalid log level, defaulting to info")
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
}

// Fields wraps logrus.Fields
type Fields logrus.Fields

// Debug logs a message at level Debug
func Debug(msg string, fields Fields) {
	if fields == nil {
		logrus.Debug(msg)
	} else {
		logrus.WithFields(logrus.Fields(fields)).Debug(msg)
	}
}

// Info logs a message at level Info
func Info(msg string, fields Fields) {
	if fields == nil {
		logrus.Info(msg)
	} else {
		logrus.WithFields(logrus.Fields(fields)).Info(msg)
	}
}

// Warn logs a message at level Warn
func Warn(msg string, fields Fields) {
	if fields == nil {
		logrus.Warn(msg)
	} else {
		logrus.WithFields(logrus.Fields(fields)).Warn(msg)
	}
}

// Error logs a message at level Error
func Error(msg string, fields Fields) {
	if fields == nil {
		logrus.Error(msg)
	} else {
		logrus.WithFields(logrus.Fields(fields)).Error(msg)
	}
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1
func Fatal(msg string, fields Fields) {
	if fields == nil {
		logrus.Fatal(msg)
	} else {
		logrus.WithFields(logrus.Fields(fields)).Fatal(msg)
	}
}

// WithField creates an entry from the standard logger and adds a field to it
func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple fields to it
func WithFields(fields Fields) *logrus.Entry {
	return logrus.WithFields(logrus.Fields(fields))
}

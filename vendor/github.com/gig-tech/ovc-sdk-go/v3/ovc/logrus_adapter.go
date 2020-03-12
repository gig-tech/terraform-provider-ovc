package ovc

import (
	"github.com/sirupsen/logrus"
)

// LogrusAdapter adapts a logrus.FieldLogger to the Logger interface
type LogrusAdapter struct {
	logrus.FieldLogger
}

// WithField implements the eponymous method of Logger
func (l LogrusAdapter) WithField(key string, value interface{}) Logger {
	return LogrusAdapter{l.FieldLogger.WithField(key, value)}
}

// WithFields implements the eponymous method of Logger
func (l LogrusAdapter) WithFields(fields map[string]interface{}) Logger {
	return LogrusAdapter{l.FieldLogger.WithFields(fields)}
}

// WithError implements the eponymous method of Logger
func (l LogrusAdapter) WithError(err error) Logger {
	return LogrusAdapter{l.FieldLogger.WithError(err)}
}

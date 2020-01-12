package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type ILogger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})
	InfoService(service string, v ...interface{})
	InfofService(service string, format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	FatalService(service string, v ...interface{})
	FatalfService(service string, format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	ErrorService(service string, v ...interface{})
	ErrorfService(service string, format string, v ...interface{})
}

type defaultLogger struct {
}

func NewDebugLogger() ILogger {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
	})
	return &defaultLogger{}
}

func (*defaultLogger) Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

func (*defaultLogger) Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

func (*defaultLogger) Warnf(format string, v ...interface{}) {
	logrus.Warnf(format, v...)
}
func (*defaultLogger) Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

func (*defaultLogger) Fatal(v ...interface{}) {
	logrus.Fatal(v...)
}

func (*defaultLogger) Fatalf(format string, v ...interface{}) {
	logrus.Fatalf(format, v...)
}

func (*defaultLogger) FatalService(service string, v ...interface{}) {
	logrus.New().WithFields(logrus.Fields{
		"service": service,
	}).Fatal(v...)
}

func (*defaultLogger) FatalfService(service string, format string, v ...interface{}) {
	logrus.New().WithFields(logrus.Fields{
		"service": service,
	}).Fatalf(format, v...)
}

func (*defaultLogger) Debug(v ...interface{}) {
	logrus.Debug(v...)
}

func (*defaultLogger) Info(v ...interface{}) {
	logrus.Info(v...)
}

func (*defaultLogger) InfoService(service string, v ...interface{}) {
	logrus.New().WithFields(logrus.Fields{
		"service": service,
	}).Info(v...)
}

func (*defaultLogger) InfofService(service string, format string, v ...interface{}) {
	logrus.New().WithFields(logrus.Fields{
		"service": service,
	}).Infof(format, v...)
}

func (*defaultLogger) Warn(v ...interface{}) {
	logrus.Warn(v...)
}
func (*defaultLogger) Error(v ...interface{}) {
	logrus.Error(v...)
}

func (*defaultLogger) ErrorService(service string, v ...interface{}) {
	logrus.New().WithFields(logrus.Fields{
		"service": service,
	}).Error(v...)
}

func (*defaultLogger) ErrorfService(service string, format string, v ...interface{}) {
	logrus.New().WithFields(logrus.Fields{
		"service": service,
	}).Errorf(format, v...)
}

package logger

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"
	"time"
)

type Logger interface {
	Setup(fluentHost string, fluentPort int, tag string) error
	Close() error
    Info(message string, fields ...interface{})
    Warn(message string, fields ...interface{})
	Error(message string, fields ...interface{})
	Debug(message string, fields ...interface{})
	Fatal(message string, fields ...interface{})
}

type RemoteLogger struct {
	fluentLogger *fluent.Fluent
	tag          string
}

func (l* RemoteLogger)  Setup(fluentHost string, fluentPort int, tag string) (error) {
	fluentLogger, err := fluent.New(fluent.Config{
		FluentPort: fluentPort,
		FluentHost: fluentHost,
	})
	if err != nil {
		return err
	}
	l.fluentLogger = fluentLogger
	l.tag = tag
	return nil
}

func (l *RemoteLogger) Close() (error) {
	err := l.fluentLogger.Close()
	if err != nil {
		return err
	}
	return nil
}

func (l *RemoteLogger) logWithFields(level logrus.Level, message string, fields ...interface{}) {
	data := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			data[fields[i].(string)] = fields[i+1]
		}
	}
	data["time"] = time.Now().Format(time.RFC3339)
	l.log(level, message, data)
}

func (l *RemoteLogger) log(level logrus.Level, message string, data map[string]interface{}) {
	data["level"] = level.String()
	data["message"] = message
	data["time"] = time.Now().Format(time.RFC3339)
	err := l.fluentLogger.PostWithTime(l.tag, time.Now(), data)
	if err != nil {
		logrus.Errorf("Failed to send log to Fluentd: %v", err)
	}
	go logrus.WithFields(data).Log(level, message)
}

func (l *RemoteLogger) Info(message string, fields ...interface{}) {
	l.logWithFields(logrus.InfoLevel, message, fields...)
}

func (l *RemoteLogger) Warn(message string, fields ...interface{}) {
	l.logWithFields(logrus.WarnLevel, message, fields...)
}

func (l *RemoteLogger) Error(message string, fields ...interface{}) {
	l.logWithFields(logrus.ErrorLevel, message, fields...)
}

func (l *RemoteLogger) Debug(message string, fields ...interface{}) {
	l.logWithFields(logrus.DebugLevel, message, fields...)
}

func (l *RemoteLogger) Fatal(message string, fields ...interface{}) {
	l.logWithFields(logrus.FatalLevel, message, fields...)
	logrus.Fatal(message)
}

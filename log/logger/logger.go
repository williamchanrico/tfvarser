package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Fields type
type Fields map[string]interface{}

// Level type
type Level int

// level of log
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Log level
const (
	DebugLevelString = "debug"
	InfoLevelString  = "info"
	WarnLevelString  = "warn"
	ErrorLevelString = "error"
	FatalLevelString = "fatal"
)

const defaultTimeFormat = time.RFC3339

// Logger struct
type Logger struct {
	logger *logrus.Logger
	config *Config
}

// Config of logger
type Config struct {
	Level      Level
	LogFile    string
	TimeFormat string
	Caller     bool
	UseColor   bool
}

// New logger
func New(config *Config) (*Logger, error) {
	if config == nil {
		config = &Config{
			Level:  InfoLevel,
			Caller: false,
		}
	}

	if config.TimeFormat == "" {
		config.TimeFormat = defaultTimeFormat
	}

	logger, err := newLogger(config)
	if err != nil {
		return nil, err
	}
	l := Logger{
		logger: logger,
		config: config,
	}
	return &l, nil
}

// DefaultLogger return default value of logger
func DefaultLogger() *Logger {
	// error is ignored because it should not throw any error
	// test must be updated if configuration is changed
	logger, err := New(&Config{
		Level:      InfoLevel,
		UseColor:   true,
		TimeFormat: defaultTimeFormat,
	})

	// only for testing purpose
	if err != nil {
		tmpLogger := logrus.New()
		tmpLogger.Fatal(err)
	}

	return logger
}

func newLogger(config *Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// set writer to file if config.LogFile is not empty
	if config.LogFile != "" {
		err := os.MkdirAll(filepath.Dir(config.LogFile), 0750)
		if err != nil && err != os.ErrExist {
			return nil, err
		}

		file, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0655)
		if err != nil {
			return nil, err
		}

		logger.SetOutput(file)
	}

	logger.SetFormatter(&logrus.TextFormatter{
		// invert the bool
		DisableColors:   config.UseColor == false,
		TimestampFormat: config.TimeFormat,
	})

	// set caller
	logger.SetReportCaller(config.Caller)

	return logger, nil
}

func setLevel(logger *logrus.Logger, level Level) {
	switch level {
	case DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	case FatalLevel:
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

// SetLevel for setting log level
func (l *Logger) SetLevel(level Level) {
	setLevel(l.logger, level)
}

// SetLevelString set level using string instead of level
func (l *Logger) SetLevelString(level string) {
	setLevel(l.logger, StringToLevel(level))
}

// SetOutput for set logger output
func (l *Logger) SetOutput(output io.Writer) {
	l.logger.SetOutput(output)
}

// StringToLevel convert string to log level
func StringToLevel(s string) Level {
	switch strings.ToLower(s) {
	case DebugLevelString:
		return DebugLevel
	case InfoLevelString:
		return InfoLevel
	case WarnLevelString:
		return WarnLevel
	case ErrorLevelString:
		return ErrorLevel
	case FatalLevelString:
		return FatalLevel
	default:
		// TODO: make this more informative when happened
		return InfoLevel
	}
}

// LevelToString convert log level to readable string
func LevelToString(l Level) string {
	switch l {
	case DebugLevel:
		return DebugLevelString
	case InfoLevel:
		return InfoLevelString
	case WarnLevel:
		return WarnLevelString
	case ErrorLevel:
		return ErrorLevelString
	case FatalLevel:
		return FatalLevelString
	default:
		return InfoLevelString
	}
}

// Debug function
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf function
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

// Debugln function
func (l *Logger) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}

// Debugw function
func (l *Logger) Debugw(message string, fields Fields) {
	l.logger.WithFields(logrus.Fields(fields)).Debugln(message)
}

// Info function
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof function
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

// Infoln function
func (l *Logger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

// Infow function
func (l *Logger) Infow(message string, fields Fields) {
	l.logger.WithFields(logrus.Fields(fields)).Infoln(message)
}

// Warn function
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf function
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

// Warnln function
func (l *Logger) Warnln(args ...interface{}) {
	l.logger.Warnln(args...)
}

// Warnw function
func (l *Logger) Warnw(message string, fields Fields) {
	l.logger.WithFields(logrus.Fields(fields)).Warnln(message)
}

// Error function
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf function
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

// Errorln function
func (l *Logger) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}

// Errorw function
func (l *Logger) Errorw(message string, fields Fields) {
	l.logger.WithFields(logrus.Fields(fields)).Errorln(message)
}

// Fatal function
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf function
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

// Fatalln function
func (l *Logger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}

// Fatalw function
func (l *Logger) Fatalw(message string, fields Fields) {
	l.logger.WithFields(logrus.Fields(fields)).Fatalln(message)
}

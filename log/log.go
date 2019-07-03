package log

import "github.com/williamchanrico/tfvarser/log/logger"

// Fields type
type Fields map[string]interface{}

var defaultLogger *logger.Logger
var debugLogger *logger.Logger

func init() {
	defaultLogger = logger.DefaultLogger()
	debugLogger, _ = logger.New(&logger.Config{Level: logger.DebugLevel})
}

// SetLevel of log
func SetLevel(level logger.Level) {
	setLevel(level)
}

// SetLevelString to set log level using string
func SetLevelString(level string) {
	setLevel(logger.StringToLevel(level))
}

// setLevel function set the log level to the desired level for defaultLogger and debugLogger
// debugLogger level can go to any level, but not with defaultLogger
// this to make sure debugLogger to be disabled when level is > debug
// and defaultLogger to not overlap with debugLogger
func setLevel(level logger.Level) {
	if level < logger.InfoLevel {
		debugLogger.SetLevel(level)
	} else {
		defaultLogger.SetLevel(level)
		debugLogger.SetLevel(level)
	}
}

// Debug function
func Debug(args ...interface{}) {
	debugLogger.Debug(args...)
}

// Debugf function
func Debugf(format string, v ...interface{}) {
	debugLogger.Debugf(format, v...)
}

// Debugw function
func Debugw(msg string, fields Fields) {
	debugLogger.Debugw(msg, logger.Fields(fields))
}

// Print function
func Print(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Println function
func Println(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Printf function
func Printf(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Printw function
func Printw(msg string, fields Fields) {
	defaultLogger.Infow(msg, logger.Fields(fields))
}

// Info function
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof function
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

// Infow function
func Infow(msg string, fields Fields) {
	defaultLogger.Infow(msg, logger.Fields(fields))
}

// Warn function
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf function
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

// Warnw function
func Warnw(msg string, fields Fields) {
	defaultLogger.Warnw(msg, logger.Fields(fields))
}

// Error function
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf function
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

// Errorw function
func Errorw(msg string, fields Fields) {
	defaultLogger.Errorw(msg, logger.Fields(fields))
}

// // Errors function to log errors package
// func Errors(err error) {
// 	defaultLogger.Errors(err)
// }

// Fatal function
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Fatalf function
func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

// Fatalw function
func Fatalw(msg string, fields Fields) {
	defaultLogger.Fatalw(msg, logger.Fields(fields))
}

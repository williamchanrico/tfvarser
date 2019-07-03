package logger

import (
	"os"
	"path"
	"testing"
)

func TestLogText(t *testing.T) {

}

func TestLogColor(t *testing.T) {

}

func TestDefaultLogger(t *testing.T) {
	l := DefaultLogger()
	l.Error("test")
}

func TestLogOutput(t *testing.T) {
	logFile := "test.log"

	conf := Config{
		LogFile: logFile,
	}

	logger, err := New(&conf)
	if err != nil {
		t.Error(err)
		return
	}

	logger.Infoln("logger")

	wdPath, err := os.Getwd()
	if err != nil {
		t.Error(err)
		return
	}

	logPath := path.Join(wdPath, logFile)

	_, err = os.Stat(logPath)
	if err != nil {
		t.Error(err)
		return
	}

	os.Remove(logPath)
}

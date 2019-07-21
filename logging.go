package heatmap

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var logInfo = log.New(ioutil.Discard, "[INFO]  ", log.Ldate|log.Ltime|log.LUTC)
var logDebug = log.New(ioutil.Discard, "[DEBUG] ", log.Ldate|log.Ltime|log.LUTC)
var logError = log.New(ioutil.Discard, "[ERROR] ", log.Ldate|log.Ltime|log.LUTC)

var loggers = []*log.Logger{
	logInfo,
	logDebug,
	logError,
}

var logLevels = []string{
	"info",
	"debug",
	"error",
}

func changeLogLevel(newLevel string) {
	newLevel = strings.ToLower(newLevel)
	newIndex := -1
	for i, l := range logLevels {
		if l == newLevel {
			newIndex = i
		}
	}

	output := io.Writer(os.Stdout)

	for i, l := range loggers {
		if i > newIndex {
			output = ioutil.Discard
		}
		l.SetOutput(output)
	}
}

package utils

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sdslabs/gasper/types"
)

// color definitions
const (
	red      = "\x1b[31m"
	green    = "\x1b[32m"
	reset    = "\x1b[0m"
	yellow   = "\x1b[33m"
	blue     = "\x1b[34m"
	magenta  = "\x1b[35m"
	cyan     = "\x1b[36m"
	lightRed = "\x1b[91m"
)

// tag definitions
const (
	ErrorTAG = 1
	InfoTAG  = 2
	DebugTAG = 3
)

var tagToStringColored = map[int]string{
	ErrorTAG: magenta + "[" + reset + red + "ERROR" + reset + magenta + "]" + reset,
	InfoTAG:  magenta + "[" + reset + blue + "INFO" + reset + magenta + "]" + reset,
	DebugTAG: magenta + "[" + reset + cyan + "DEBUG" + reset + magenta + "]" + reset,
}

var tagToString = map[int]string{
	ErrorTAG: "[ERROR]",
	InfoTAG:  "[INFO]",
	DebugTAG: "[DEBUG]",
}

var (
	logfile, _ = os.Create("gasper.log")
	mutex      = sync.Mutex{}
)

func coloredContext(context string) string {
	return magenta + "(" + reset + cyan + context + reset + magenta + ")" + reset
}

func out(context, message string, TAG int) {
	currentTime := time.Now()
	hours := fmt.Sprintf("%d", currentTime.Hour())
	if currentTime.Hour() < 10 {
		hours = "0" + hours
	}
	minutes := fmt.Sprintf("%d", currentTime.Minute())
	if currentTime.Minute() < 10 {
		minutes = "0" + minutes
	}
	seconds := fmt.Sprintf("%d", currentTime.Second())
	if currentTime.Second() < 10 {
		seconds = "0" + seconds
	}
	timeLog := fmt.Sprintf(
		"%d-%d-%d %s:%s:%s",
		currentTime.Day(),
		currentTime.Month(),
		currentTime.Year(),
		hours,
		minutes,
		seconds,
	)
	fmt.Println(tagToStringColored[TAG] + " " + coloredContext(context) + " " + yellow + timeLog + reset + lightRed + " >>> " + reset + green + message + reset)
	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		logfile.WriteString(fmt.Sprintf("%s (%s) %s >>> %s\n", tagToString[TAG], context, timeLog, message))
	}()
}

// Log logs to the console with your custom TAG
func Log(context, message string, TAG int) {
	out(context, message, TAG)
}

// LogInfo logs information to the console
func LogInfo(context, template string, args ...interface{}) {
	out(context, fmt.Sprintf(template, args...), InfoTAG)
}

// LogDebug logs debug messages to console
func LogDebug(context, template string, args ...interface{}) {
	out(context, fmt.Sprintf(template, args...), DebugTAG)
}

// LogError logs type error to console
func LogError(context string, err error) {
	if err == nil {
		return
	}
	out(context, err.Error(), ErrorTAG)
}

// LogResErr logs type ResponseError to console
func LogResErr(context string, e types.ResponseError) {
	message := fmt.Sprintf("%d: %s\n%s", e.Status(), e.Message(), e.Verbose())
	out(context, message, ErrorTAG)
}

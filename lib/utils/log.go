package utils

import (
	"fmt"
	"time"

	"github.com/sdslabs/SWS/lib/types"
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
	ErrorTAG = magenta + "[" + reset + red + "ERROR" + reset + magenta + "]"
	InfoTAG  = magenta + "[" + reset + blue + "INFO" + reset + magenta + "]"
	DebugTAG = magenta + "[" + reset + cyan + "DEBUG" + reset + magenta + "]"
)

func out(s, TAG string) {
	currentTime := time.Now()
	timeLog := fmt.Sprintf(
		"%d-%d-%d %d:%d:%d",
		currentTime.Day(),
		currentTime.Month(),
		currentTime.Year(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second(),
	)
	fmt.Println(TAG + reset + " " + yellow + timeLog + reset + lightRed + " >>> " + reset + green + s + reset)
}

// Log logs to the console with your custom TAG
func Log(s, TAG string) {
	out(s, TAG)
}

// LogInfo logs information to the console
func LogInfo(f string, v ...interface{}) {
	s := fmt.Sprintf(f, v...)
	out(s, InfoTAG)
}

// LogDebug logs debug messages to console
func LogDebug(f string, v ...interface{}) {
	s := fmt.Sprintf(f, v...)
	out(s, DebugTAG)
}

// LogError logs type error to console
func LogError(e error) {
	s := e.Error()
	out(s, ErrorTAG)
}

// LogResErr logs type ResponseError to console
func LogResErr(e types.ResponseError) {
	s := fmt.Sprintf("%d: %s\n%s", e.Status(), e.Message(), e.Verbose())
	out(s, ErrorTAG)
}

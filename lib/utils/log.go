package utils

import (
	"fmt"
	"time"

	"github.com/sdslabs/SWS/lib/configs"
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

func out(s, tag string) {
	if configs.SWSConfig["debug"].(bool) {
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
		fmt.Println(tag + reset + " " + yellow + timeLog + reset + lightRed + " >>> " + reset + green + s + reset)
	}
}

// Log string to the console
func Log(s string) {
	out(s, InfoTAG)
}

// Logf is Log with format string
func Logf(f string, v ...interface{}) {
	s := fmt.Sprintf(f, v...)
	out(s, InfoTAG)
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

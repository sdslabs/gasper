package utils

import (
	"fmt"

	"github.com/sdslabs/SWS/lib/configs"
	"github.com/sdslabs/SWS/lib/types"
)

const (
	red   = "\x1b[31m"
	green = "\x1b[32m"
	reset = "\x1b[0m"
)

func out(s string) {
	if configs.SWSConfig["debug"].(bool) {
		fmt.Println(red + ">>> " + reset + green + s + reset)
	}
}

// Log string to the console
func Log(s string) {
	out(s)
}

// Logf is Log with format string
func Logf(f string, v ...interface{}) {
	s := fmt.Sprintf(f, v...)
	out(s)
}

// LogError logs type error to console
func LogError(e error) {
	s := e.Error()
	out(s)
}

// LogResErr logs type ResponseError to console
func LogResErr(e types.ResponseError) {
	s := fmt.Sprintf("%d: %s\n%s", e.Status(), e.Message(), e.Verbose())
	out(s)
}

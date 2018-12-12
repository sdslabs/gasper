package utils

import (
	"fmt"
	"strings"
)

// NewExecCmd creates new command to be executed inside container
func NewExecCmd(commands []string) []string {
	return []string{"bash", "-c", strings.Join(commands, ";")}
}

// NewExecCmdInApp is same as NewExecCmd with 'cd app'
func NewExecCmdInApp(commands []string, appDir string) []string {
	cd := []string{fmt.Sprintf("cd %s", appDir)}
	cmd := append(cd, commands...)
	return NewExecCmd(cmd)
}

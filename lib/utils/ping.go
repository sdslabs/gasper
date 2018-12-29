package utils

import (
	"fmt"
	"net/http"
)

// NotAlive checks if a given instance is alive or not
func NotAlive(url string) bool {
	_, err := http.Get(url)
	fmt.Println(err)
	if err != nil {
		return true
	}
	return false
}

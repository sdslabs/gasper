package utils

import (
	"fmt"
	"net"
	"time"
)

// NotAlive checks if a given instance is alive or not
func NotAlive(url string) bool {
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", url)
	if err != nil {
		fmt.Println(err)
		return true
	}
	conn.Close()
	return false
}

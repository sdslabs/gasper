package utils

import (
	"net"
)

// GetOutboundIP returns the preferred outbound IP of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// HostIP variable stores the IPv4 address of the host machine
var HostIP = GetOutboundIP()

package utils

import (
	"net"

	"github.com/sdslabs/SWS/lib/configs"
)

// GetOutboundIP returns the preferred outbound IP of this machine
func GetOutboundIP() string {
	if configs.SWSConfig["offlineMode"].(bool) {
		return "0.0.0.0"
	}
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic("The machine is not connected to any network")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// HostIP variable stores the IPv4 address of the host machine
var HostIP = GetOutboundIP()

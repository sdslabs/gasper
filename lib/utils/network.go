package utils

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/types"
)

// HostIP variable stores the IPv4 address of the host machine
var HostIP, _ = GetOutboundIP()

// GetOutboundIP returns the preferred outbound IP of this machine
func GetOutboundIP() (string, error) {
	if configs.GasperConfig.OfflineMode {
		return "0.0.0.0", nil
	}
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		LogInfo("Utils-Network-1", "This node is not connected to any network")
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

// NotAlive checks if a given instance is alive or not
func NotAlive(url string) bool {
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", url)
	if err != nil {
		Log("Utils-Network-2", "Health-Check failed for instance "+url, ErrorTAG)
		LogError("Utils-Network-3", err)
		return true
	}
	defer conn.Close()
	return false
}

// IsGenDNSAlive checks if a GenDNS instance is alive or not
func IsGenDNSAlive(url string) bool {
	target := fmt.Sprintf("%s.%s", types.Master, configs.GasperConfig.Domain)
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)
	if _, _, err := c.Exchange(&m, url); err != nil {
		Log("Utils-Network-4", "Health-Check failed for GenDNS instance "+url, ErrorTAG)
		LogError("Utils-Network-5", err)
		return false
	}
	return true
}

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetFreePorts asks the kernel for free open ports that are ready to use.
func GetFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}

// IsValidPort checks if the port is valid and free to use
func IsValidPort(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}

	err = ln.Close()
	if err != nil {
		return false
	}

	return true
}

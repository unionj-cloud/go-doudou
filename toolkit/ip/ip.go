package ip

import (
	"net"
)

// GetOutboundIP return local public net.IP
func GetOutboundIP() net.IP {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

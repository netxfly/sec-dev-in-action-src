package scanner

import (
	"fmt"
	"net"
	"time"
)

func Connect(ip string, port int) (string, int, error) {

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 2*time.Second)

	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	return ip, port, err
}

package scanner

import (
	"fmt"
	"net"
	"time"
)

func Connect(ip string, port int) (net.Conn, error) {

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 1*time.Second)

	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	return conn, err
}

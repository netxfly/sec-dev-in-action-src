package plugins

import (
	"fmt"
	"net"

	"sec-dev-in-action-src/scanner/password_crack/models"
	"sec-dev-in-action-src/scanner/password_crack/vars"

	"golang.org/x/crypto/ssh"
)

func ScanSsh(s models.Service) (result models.ScanResult, err error) {
	result.Service = s
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		Timeout: vars.TimeOut,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port), config)
	if err != nil {
		return result, err
	}

	session, err := client.NewSession()
	if err != nil {
		return result, err

	}
	err = session.Run("echo 666")
	if err != nil {
		return result, err
	}

	result.Result = true

	defer func() {
		if client != nil {
			_ = client.Close()
		}
		if session != nil {
			_ = session.Close()
		}
	}()

	return result, err
}

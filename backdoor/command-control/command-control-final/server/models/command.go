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

package models

import "time"

type (
	Command struct {
		Id         int64     `json:"id"`
		AgentId    string    `json:"agent_id"`
		Content    string    `json:"content"`
		Status     int       `json:"status"`
		CreateTime time.Time `xorm:"created"`
		UpdateTime time.Time `xorm:"updated"`
		Result     string    `json:"result"`
	}
)

func NewCommand(agentId string, content string) *Command {
	return &Command{
		Id:      0,
		AgentId: agentId,
		Content: content,
		Status:  0,
	}
}

func (c *Command) Insert() error {
	_, err := Engine.Insert(c)
	return err
}

func ListCommandByAgentId(agentId string) ([]Command, error) {
	cmds := make([]Command, 0)
	err := Engine.Where("agent_id=? and status=0", agentId).Find(&cmds)
	return cmds, err
}

func ListFinishCommand() ([]Command, error) {
	cmds := make([]Command, 0)
	err := Engine.Where("status=1").Find(&cmds)
	return cmds, err
}

func UpdateCommandResult(cmdId int64, result string) error {
	cmd := new(Command)
	has, err := Engine.Id(cmdId).Get(cmd)
	if err != nil {
		return err
	}
	if has {
		cmd.Result = result
		_, err = Engine.Id(cmdId).Update(cmd)
	}

	return err
}

func UpdateCommandStatus(cmdId int64, status int) error {
	cmd := new(Command)
	has, err := Engine.Id(cmdId).Get(cmd)
	if err != nil {
		return err
	}
	if has {
		cmd.Status = status
		_, err = Engine.Id(cmdId).Update(cmd)
	}

	return err
}

// 命令执行完的标志
func SetCmdStatusToFinished(cmdId int64) error {
	err := UpdateCommandStatus(cmdId, 1)
	return err
}

// 在控制台展示完状态的标志
func SetCmdStatusToEnd(cmdId int64) error {
	err := UpdateCommandStatus(cmdId, 2)
	return err
}

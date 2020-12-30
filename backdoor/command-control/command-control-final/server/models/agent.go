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

type Agent struct {
	Id           int64
	AgentId      string    `json:"agent_id"`
	Platform     string    `json:"platform"`
	Architecture string    `json:"architecture"`
	UserName     string    `json:"user_name"`
	UserGUID     string    `json:"user_guid"`
	HostName     string    `json:"host_name"`
	Ips          []string  `json:"ips" xorm:"text"`
	Pid          int       `json:"pid"`
	Debug        bool      `json:"debug"`
	Proto        string    `json:"proto"`
	UserAgent    string    `json:"user_agent"`
	Initial      bool      `json:"initial"`
	CreateTime   time.Time `xorm:"created"`
	UpdateTime   time.Time `xorm:"updated"`
	Version      int       `xorm:"version"`
}

func (a *Agent) Insert() error {
	_, err := Engine.Insert(a)
	return err
}

func ListAgents() ([]Agent, error) {
	agents := make([]Agent, 0)
	err := Engine.Find(&agents)
	return agents, err
}

func UpdateAgent(agentId string) error {
	agent := new(Agent)
	has, err := Engine.Where("agent_id=?", agentId).Get(agent)
	if err != nil {
		return err
	}
	if has {
		_, err = Engine.Id(agent.Id).Update(agent)
	}
	return err
}

func ExistAgentId(agentId string) (bool, error) {
	agent := new(Agent)
	has, err := Engine.Where("agent_id=?", agentId).Get(agent)
	return has, err
}

func RemoveAll() error {
	_, err := Engine.Exec("delete from agent")
	return err
}

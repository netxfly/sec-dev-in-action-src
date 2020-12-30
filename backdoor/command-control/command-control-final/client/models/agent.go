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

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"

	"golang.org/x/net/http2"

	uuid "github.com/satori/go.uuid"
)

type Agent struct {
	AgentId      uuid.UUID    `json:"agent_id"`
	Platform     string       `json:"platform"`
	Architecture string       `json:"architecture"`
	UserName     string       `json:"user_name"`
	UserGUID     string       `json:"user_guid"`
	HostName     string       `json:"host_name"`
	Ips          []string     `json:"ips"`
	Pid          int          `json:"pid"`
	Debug        bool         `json:"debug"`
	Proto        string       `json:"proto"`
	Client       *http.Client `json:"client"`
	UserAgent    string       `json:"user_agent"`
	Initial      bool         `json:"initial"`
	URL          string       `json:"url"`
	Host         string       `json:"host"`
}

type AgentInfo struct {
	Id           int64
	AgentId      uuid.UUID `json:"agent_id"`
	Platform     string    `json:"platform"`
	Architecture string    `json:"architecture"`
	UserName     string    `json:"user_name"`
	UserGUID     string    `json:"user_guid"`
	HostName     string    `json:"host_name"`
	Ips          []string  `json:"ips"`
	Pid          int       `json:"pid"`
	Debug        bool      `json:"debug"`
	Proto        string    `json:"proto"`
	UserAgent    string    `json:"user_agent"`
	Initial      bool      `json:"initial"`
}

func NewAgent(debug bool, protocol string) (*Agent, error) {
	uuidV4, _ := uuid.NewV1()

	agent := &Agent{
		AgentId:      uuidV4,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		Ips:          nil,
		Pid:          os.Getpid(),
		Debug:        debug,
		Proto:        protocol,
		Client:       nil,
		UserAgent:    "Mozilla / 5.0(Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.25 Safari/537.36",
		Initial:      false,
		URL:          "",
		Host:         "",
	}

	switch protocol {
	case "h2c":
		agent.URL = fmt.Sprintf("%v://127.0.0.1:8080", "https")
	case "http2":
		agent.URL = fmt.Sprintf("%v://127.0.0.1:8080", "https")
	case "http":
		agent.URL = fmt.Sprintf("%v://127.0.0.1:8080", "http")
	}

	u, err := user.Current()
	if err != nil {
		return agent, err
	}

	agent.UserName = u.Username
	agent.UserGUID = u.Gid

	h, errH := os.Hostname()
	if errH != nil {
		return agent, err
	}
	agent.HostName = h

	interfaces, err := net.Interfaces()
	if err != nil {
		return agent, err
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				if IsIPv4(addr.String()) {
					agent.Ips = append(agent.Ips, addr.String())
				}
			}
		} else {
			return agent, err
		}
	}

	agent.Client = getClient(protocol)

	return agent, err
}

func (a *Agent) ParseInfo() AgentInfo {
	return AgentInfo{
		Id:           0,
		AgentId:      a.AgentId,
		Platform:     a.Platform,
		Architecture: a.Architecture,
		UserName:     a.UserAgent,
		UserGUID:     a.UserGUID,
		HostName:     a.HostName,
		Ips:          a.Ips,
		Pid:          a.Pid,
		Debug:        a.Debug,
		Proto:        a.Proto,
		UserAgent:    a.UserAgent,
		Initial:      a.Initial,
	}
}

func getClient(protocol string) *http.Client {
	client := &http.Client{}
	switch protocol {
	case "http":
	case "h2c":

	case "http2":
		client = &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	return client
}

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

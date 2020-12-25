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

package services

import (
	"sync"
	"time"

	redis_server "sec-dev-in-action-src/honeypot/server/services/redis-server"

	"sec-dev-in-action-src/honeypot/server/services/mysql-server"
	"sec-dev-in-action-src/honeypot/server/services/ssh-server"
	"sec-dev-in-action-src/honeypot/server/services/web"
	"sec-dev-in-action-src/honeypot/server/vars"
)

type
(
	HoneypotServices func(string, bool) error

	ServiceInfo struct {
		ServerName       string           `json:"server_name"`
		ListenAddr       string           `json:"listen_addr"`
		Flag             bool             `json:"flag"`
		HoneypotServices HoneypotServices `json:"honeypot_services"`
	}
)

var (
	Services   []ServiceInfo
	fnServices = map[string]HoneypotServices{"ssh": ssh_server.StartSsh, "redis": redis_server.StartRedis, "mysql": mysql.StartMysql, "web": web.StartWeb}
)

func init() {
	Services = make([]ServiceInfo, 0)
	for service, item := range vars.Config.Services {
		Services = append(Services, ServiceInfo{ServerName: service, ListenAddr: item.Addr, HoneypotServices: fnServices[service], Flag: item.Flag})
	}
}

func Start() {
	var wg sync.WaitGroup
	for _, s := range Services {
		wg.Add(1)
		go func(service ServiceInfo) {
			err := service.HoneypotServices(service.ListenAddr, service.Flag)
			_ = err
			wg.Done()
		}(s)
	}
	wg.Wait()

	for {
		time.Sleep(100 * time.Second)
	}
}

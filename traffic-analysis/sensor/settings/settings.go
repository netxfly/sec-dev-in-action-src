/*

Copyright (c) 2017 xsec.io

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

package settings

import (
	"gopkg.in/ini.v1"

	"sec-dev-in-action-src/traffic-analysis/sensor/misc"
)

var (
	Cfg        *ini.File
	DeviceName string
	DebugMode  bool
	FilterRule string
	Ips        []string
)

func init() {
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)

	if err != nil {
		misc.Log.Panicln(err)
	}

	DeviceName = Cfg.Section("").Key("DEVICE_NAME").MustString("eth0")
	DebugMode = Cfg.Section("").Key("DEBUG_MODE").MustBool(false)
	FilterRule = Cfg.Section("").Key("FILTER_RULE").MustString("tcp or (udp and dst port 53)")

	Ips, _ = misc.GetIpList(DeviceName)
	misc.Log.Infof("Device name:[%v], ip addr:%v, Debug mode:[%v]", DeviceName, Ips, DebugMode)
}

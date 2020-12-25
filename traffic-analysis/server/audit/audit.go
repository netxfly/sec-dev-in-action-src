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

package audit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"sec-dev-in-action-src/traffic-analysis/server/models"
	"sec-dev-in-action-src/traffic-analysis/server/settings"
	"sec-dev-in-action-src/traffic-analysis/server/util"
)

var (
	EvilIpUrl string
)

func init() {
	sec := settings.Cfg.Section("EVIL_IPS")
	EvilIpUrl = sec.Key("API_URL").MustString("")
}

func PacketAduit(sensorIp string, connInfo models.ConnectionInfo) (err error, result bool, detail models.IplistApi) {
	ips := make([]string, 0)
	ips = append(ips, connInfo.SrcIp, connInfo.DstIp)

	for _, ip := range ips {
		if ip == sensorIp {
			continue
		}
		evilUrl := fmt.Sprintf("%v/api/ip/%v", EvilIpUrl, ip)
		resp, err := http.Get(evilUrl)
		var detail models.IplistApi
		if err == nil {
			ret, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				err = json.Unmarshal(ret, &detail)
				result = detail.Evil
				// util.Log.Debugf("check ip:%v, result: %v, detail: %v", ip, result, detail)
				if result {
					evilConnInfo := models.NewEvilConnectionInfo(sensorIp, connInfo, detail)
					evilConnInfo.Insert()
				}
			}

		}
	}

	return err, result, detail
}

func HttpAudit(sensorIp string, req models.HttpReq) (err error, result bool, evilReq *models.EvilHttpReq) {
	util.Log.Debugf("sensorIp: %v, req: %v", sensorIp, req)
	evilReq = models.NewEvilHttpReq(sensorIp, result, req)
	evilReq.Insert()
	return err, result, evilReq
}

func DnsAudit(sensorIp string, dns models.Dns) (err error, result bool, evilDns *models.EvilDns) {
	util.Log.Debugf("sensorIp: %v, req: %v", sensorIp, dns)
	evilDns = models.NewEvilDns(sensorIp, result, dns)
	err = evilDns.Insert()
	return err, result, evilDns
}

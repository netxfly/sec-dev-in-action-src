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

package api

import (
	"sec-dev-in-action-src/proxy-honeypot/agent/settings"

	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	SECRET string
	APIURL string
)

// md5 function
func MD5(s string) (m string) {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// create a sign by key & md5
func MakeSign(t string, key string) (sign string) {
	sign = MD5(fmt.Sprintf("%s%s", t, key))
	return sign
}

func init() {
	cfg := settings.Cfg
	sec := cfg.Section("server")
	SECRET = sec.Key("SECRET").MustString("SECRET")
	APIURL = sec.Key("API_URL").MustString("http://127.0.0.1/api/send")
}

func Post(data string) (err error) {
	t := time.Now().Format("2006-01-02 15:04:05")
	hostName, _ := os.Hostname()
	_, err = http.PostForm(APIURL, url.Values{"timestamp": {t}, "secureKey": {MakeSign(t, SECRET)}, "data": {data}, "hostname": {hostName}})
	return err
}

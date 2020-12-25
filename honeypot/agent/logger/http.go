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

package logger

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	"sec-dev-in-action-src/honeypot/agent/settings"
	"sec-dev-in-action-src/honeypot/agent/util"
)

type (
	HttpHook struct {
		HttpClient http.Client
	}
)

func NewHttpHook() (*HttpHook, error) {
	timeout := 1 * time.Second
	client := http.Client{Timeout: timeout}

	return &HttpHook{HttpClient: client}, nil
}

func (hook *HttpHook) Fire(entry *logrus.Entry) (err error) {
	var serverUrl string
	field := entry.Data
	serverUrl = settings.ManagerUrl

	_, ok := field["api"]
	if ok {
		urlApi := fmt.Sprintf("%v%v", serverUrl, field["api"])
		data := entry.Message
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		secureKey := util.MakeSign(time.Now().Format("2006-01-02 15:04:05"), settings.SecKey)
		resp, err := hook.HttpClient.PostForm(urlApi, url.Values{"timestamp": {timestamp}, "secureKey": {secureKey}, "data": {data}})
		fmt.Printf("resp: %v, err: %v\n", resp, err)

	}

	return err
}

func (hook *HttpHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

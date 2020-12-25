package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"exchange_zero_trust_api/logger"
	"exchange_zero_trust_api/vars"
)

func SendSMS(username, phone, deviceId, code, content, srcIp string) (result bool, err error) {
	client := &http.Client{}
	postData := strings.NewReader(fmt.Sprintf("recipients=%s&content=%s", phone, content))
	req, err := http.NewRequest("POST", vars.SmsURL, postData)
	if err == nil {
		req.Header.Add(vars.SmsHeader, vars.SmsKey)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err == nil {
			result = true
			defer resp.Body.Close()
			ret, err := ioutil.ReadAll(resp.Body)
			logger.Log.Infof("user:%v, deviceId: %v, code: %v, sms send result: %v, err: %v",
				username, deviceId, code, strings.TrimSpace(string(ret)), err)
		}
	}
	return result, err
}

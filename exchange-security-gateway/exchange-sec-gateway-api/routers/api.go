package routers

import (
	"encoding/json"
	"github.com/labstack/echo"

	"io/ioutil"

	"exchange_zero_trust_api/logger"
	"exchange_zero_trust_api/util"
	"exchange_zero_trust_api/vars"
)

// wbxml协议解码
func DecodeWbxml(c echo.Context) error {
	d, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		err = c.JSON(500, err.Error())
	}

	deviceInfo, err := util.Parse(string(d))
	info, err := json.Marshal(deviceInfo)
	_ = err
	return c.String(200, string(info))
}

// 发送短信与邮件通知
func SendNotice(c echo.Context) error {
	username := c.FormValue("username")
	phone := c.FormValue("phone")
	content := c.FormValue("content")
	deviceId := c.FormValue("device_id")
	code := c.FormValue("code")
	srcIp := c.FormValue("src_ip")

	logger.Log.Warnf("code: %v, type: %v, user: %v, content: %v",
		code,
		vars.SendType,
		username,
		content,
	)

	var (
		result bool
		err    error
	)

	switch vars.SendType {
	case "sms":
		result, err = util.SendSMS(username, phone, deviceId, code, content, srcIp)
	case "mail":
		result, err = util.SendMail(username, content)
	case "neixin":
		// result, err = util.SendNX(username, content)
	case "weixin":
		result, err = util.SendWeiXin(username, content)
	case "dingding":
		result, err = util.SendDingDingMessage(username, content)
	default:
		result, err = util.SendSMS(username, phone, deviceId, code, content, srcIp)
	}

	if err == nil && result {
		return c.String(200, "ok")
	} else {
		return c.String(500, "error")
	}
}

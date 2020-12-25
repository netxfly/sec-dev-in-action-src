package routers

import (
	"encoding/json"
	"exchange_zero_trust_api/logger"
	"exchange_zero_trust_api/models"
	"fmt"

	"net/http"

	"github.com/labstack/echo"
)

func EwsActivation(c echo.Context) error {
	code := c.Param("code")
	var (
		username   string
		state      int
		ip         string
		ipList     string
		clientType string
		ipStatus   bool

		ewsResponse   models.EwsResponse
		ewsActiveResp models.EwsActiveResp
	)

	exist, ewsCode, err := models.ExistEwsActiveCode(code)
	logger.Log.Infof("exist: %v, Active Code: %v, err: %v", exist, ewsCode, err)

	if !exist {
		ewsResponse.ClientType = clientType
		ewsResponse.ActiveStatus = "STATE_INVALID"
		err = c.Render(http.StatusOK, "ews_state.html", ewsResponse)
	} else {
		// 激活码中的信息获取
		{
			username = ewsCode.Username
			ip = ewsCode.Ip
			clientType = ewsCode.ClientType
			state = ewsCode.State
			ipList = ewsCode.IpList
		}

		// 激活状态页面的渲染数据
		{
			ewsResponse.Username = username
			ewsResponse.ClientType = clientType
			ewsResponse.ClientIp = ip
			ewsResponse.IpList = ipList
		}

		// 激活页面的渲染数据
		{
			ewsActiveResp.Username = username
			ewsActiveResp.ClientType = clientType
			ewsActiveResp.ClientIp = ip
			ewsActiveResp.IpList = ipList
			ewsActiveResp.Code = code
		}

		ipStatus, err = models.GetEwsIpStatus(username, ip)
		logger.Log.Errorf("code status: %v, ip status: %v, err: %v",
			state,
			ipStatus,
			err,
		)

		switch state {
		case 2:
			if ipStatus {
				ewsResponse.ActiveStatus = "STATE_ACTIVED"
				err = c.Render(http.StatusOK, "ews_state.html", ewsResponse)
			} else {
				ewsResponse.ActiveStatus = "STATE_REJECTED"
				err = c.Render(http.StatusOK, "ews_state.html", ewsResponse)
			}
		case 0:
			//激活码尚未使用 or 设备地址有效的情况下，显示激活页面
			if ipStatus {
				err = c.Render(http.StatusOK, "active_ews.html", ewsActiveResp)
			} else {
				ewsResponse.ActiveStatus = "STATE_INVALID"
				err = c.Render(http.StatusOK, "ews_state.html", ewsResponse)
			}
		default:
			ewsResponse.ActiveStatus = "STATE_INVALID"
			err = c.Render(http.StatusOK, "ews_state.html", ewsResponse)
		}
	}

	return err
}

func ActiveEws(c echo.Context) error {
	code := c.FormValue("c")
	respData := RespData{Code: 500, Data: "", Message: "激活码不存在"}
	logger.Log.Infof("code: %v", code)

	if code != "" {
		result, ewsCode, err := models.VerifyEwsActiveCode(code)
		logger.Log.Infof("result: %v, ewsCode: %v, err: %v", result, ewsCode, err)
		if result && ewsCode.State == 0 {
			username := ewsCode.Username
			ip := ewsCode.Ip
			clientType := ewsCode.ClientType
			// 设置激活码为已使用
			err := models.SetEwsActiveCodeState(code, username, ip)
			// 添加EWS可信IP
			// ews_type=0, 表示自动加入，ews_type=1，表示手动激活加入
			// state=0, 表示激活，state=1, 表示未激活，state=-1, 表示禁用
			err = models.AddEwsTrustIp(username, ip, 1, 0, clientType)
			_ = err
			respData.Code = 200
			respData.Message = fmt.Sprintf("已经允许来自%v的邮件客户端%v访问邮件", ip, clientType)
		}
	}

	data, err := json.Marshal(respData)
	err = c.String(http.StatusOK, string(data))

	return err
}

func IgnoreEws(c echo.Context) error {
	code := c.FormValue("c")
	respData := RespData{Code: 500, Data: "", Message: "激活码不存在"}
	logger.Log.Infof("code: %v", code)

	if code != "" {
		result, ewsCode, err := models.VerifyEwsActiveCode(code)
		logger.Log.Infof("result: %v, ewsCode: %v, err: %v", result, ewsCode, err)
		if result && ewsCode.State == 0 {
			username := ewsCode.Username
			ip := ewsCode.Ip
			clientType := ewsCode.ClientType
			// 设置激活码为已使用
			err := models.SetEwsActiveCodeState(code, username, ip)
			// 添加EWS可信IP
			// ews_type=0, 表示自动加入，ews_type=1，表示手动激活加入
			// state=0, 表示激活，state=1, 表示未激活，state=-1, 表示禁止
			err = models.AddEwsTrustIp(username, ip, 1, -1, clientType)
			_ = err
			respData.Code = 200
			respData.Message = fmt.Sprintf("已经禁止来自%v的邮件客户端%v访问邮件", ip, clientType)
		}
	}

	data, err := json.Marshal(respData)
	err = c.String(http.StatusOK, string(data))

	return err
}

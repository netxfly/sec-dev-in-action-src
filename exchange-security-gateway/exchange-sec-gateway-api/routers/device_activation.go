package routers

import (
	"exchange_zero_trust_api/logger"
	"exchange_zero_trust_api/models"

	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

type (
	ApiData struct {
		Name         string
		DeviceModel  string
		DeviceType   string
		DeviceId     string
		Imei         string
		PhoneNumber  string
		ActiveStatus string
		Code         string
	}

	RespData struct {
		Code    int
		Data    string
		Message string
	}
)

func Activation(c echo.Context) error {
	code := c.Param("code")
	var (
		deviceName string
		deviceMode string
		imei       string
		deviceId   string
		deviceNum  int

		apiDate ApiData
		err     error
	)

	if code != "" {
		// 检查激活码
		has, user, deviceInfo, err := models.CheckActiveCode(code)
		// logger.Log.Infof("has: %v, user: %v, device_info: %v, err: %v", has, user, deviceInfo, err)
		// 激活码是存在的
		if err == nil && has {
			// 获取当前用户的设备数
			deviceNum = models.GetDeviceNum(user)
			deviceId = deviceInfo.DeviceId

			apiDate.Name = user
			apiDate.DeviceId = deviceId
			apiDate.DeviceModel = deviceInfo.DeviceType
			apiDate.DeviceType = deviceInfo.DeviceType
			apiDate.Code = code

			// 取出手机端通过WBXML协议获取的数据
			wbxmlInfo, err := models.GetDeviceInfoByDeviceId(deviceId)
			if err == nil {
				deviceName = wbxmlInfo.FriendlyName
				deviceName = wbxmlInfo.Model
				imei = wbxmlInfo.IMEI
				phone := wbxmlInfo.PhoneNumber

				apiDate.DeviceModel = deviceMode
				apiDate.Imei = imei
				apiDate.DeviceType = deviceName
				apiDate.PhoneNumber = phone
			}
			CodeStatus := deviceInfo.State
			deviceStatue := models.GetDeviceState(user, deviceId)
			// logger.Log.Infof("user: %v, deviceId: %v, CodeStatus: %v, deviceStatue: %v",
			//	user, deviceId, CodeStatus, deviceStatue)

			switch CodeStatus {
			case 0:
				// 显示已激活的状态页面
				if deviceStatue == 0 {
					apiDate.ActiveStatus = "STATE_ACTIVED"
					err = c.Render(http.StatusOK, "deviceState.html", apiDate)
					// 显示已拒绝状态的页面
				} else if deviceStatue == 3 {
					apiDate.ActiveStatus = "STATE_REJECTED"
					err = c.Render(http.StatusOK, "deviceState.html", apiDate)
				}
			case 3:
				if deviceStatue == 3 {
					apiDate.ActiveStatus = "STATE_REJECTED"
					err = c.Render(http.StatusOK, "deviceState.html", apiDate)
				}
			default:
				// 激活的设备数已经超过10个
				if deviceNum >= 10 {
					apiDate.ActiveStatus = "STATE_EXCEED"
					err = c.Render(http.StatusOK, "deviceState.html", apiDate)
				} else {
					// 显示激活页面
					err = c.Render(http.StatusOK, "activeSync.html", apiDate)
				}
			}
		} else {
			// 激活码不存在
			apiDate.ActiveStatus = "STATE_INVALID"
			err = c.Render(http.StatusOK, "deviceState.html", apiDate)
		}
	} else {
		//	激活码为空
		apiDate.ActiveStatus = "STATE_INVALID"
		_ = c.Render(http.StatusOK, "deviceState", apiDate)
	}

	return err
}

func ActiveDevice(c echo.Context) error {
	code := c.FormValue("c")
	respData := RespData{Code: 500, Data: "", Message: "设备不存在"}
	logger.Log.Infof("code: %v", code)
	if code != "" {
		result, user, device := models.VerifyActiveCode(code)
		logger.Log.Infof("result: %v, user: %v, device: %v", result, user, device)
		if result {
			err := models.ActiveDevice(user, device.DeviceId)
			result, err := models.ResetActiveCodeStatus(code, user, device, 0)
			// logger.Log.Infof("result: %v, err: %v", result, err)
			_ = err
			_ = result
			respData.Code = 200
			respData.Message = "已经允许设备访问"
		}
	}

	data, err := json.Marshal(respData)
	err = c.String(http.StatusOK, string(data))

	return err
}

func IgnoreDevice(c echo.Context) error {
	code := c.FormValue("c")
	respData := RespData{Code: 500, Data: "", Message: "设备不存在"}

	if code != "" {
		result, user, device := models.VerifyActiveCode(code)
		if result {
			err := models.IgnoreDevice(user, device.DeviceId)
			result, err := models.ResetActiveCodeStatus(code, user, device, 3)
			_ = err
			_ = result
			respData.Code = 200
			respData.Message = "已忽略该设备"
		}
	}

	data, err := json.Marshal(respData)
	err = c.String(http.StatusOK, string(data))
	return err
}

func NotFound(c echo.Context) error {
	err := c.NoContent(http.StatusNotFound)
	return err
}

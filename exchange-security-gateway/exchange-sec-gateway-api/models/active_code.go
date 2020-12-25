package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"exchange_zero_trust_api/logger"
	"exchange_zero_trust_api/vars"
)

// 删除激活码
func RemoveActiveCode(code string) (err error) {
	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.CodePrefix, code)
	InitRedis()
	_, err = vars.RedisInstance.Del(key).Result()
	return err
}

// 设置激活码的值及超时时间，该函数不对外公开
func setActiveCode(device Device, code string) (err error) {
	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.CodePrefix, code)
	deviceStr, err := json.Marshal(device)
	v := fmt.Sprintf("%v-_-%v", device.User, string(deviceStr))
	InitRedis()
	_, err = vars.RedisInstance.Set(key, v, 60*60*12*time.Second).Result()
	return err
}

// 获取激活码的值
func getActiveCodeValue(code string) (has bool, user string, deviceInfo Device, err error) {
	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.CodePrefix, code)
	InitRedis()

	c, err := vars.RedisInstance.Get(key).Result()
	logger.Log.Infof("key: %v, c: %v, err: %v", key, c, err)
	if err == nil {
		t := strings.Split(c, "-_-")
		if len(t) == 2 {
			user = t[0]
			deviceStr := t[1]
			err = json.Unmarshal([]byte(deviceStr), &deviceInfo)
			logger.Log.Infof("user: %v, deviceStr: %v, device: %v, err: %v",
				user, deviceStr, deviceInfo, err)
			if err == nil {
				deviceInfo.User = user
				has = true
			}
		}
	}
	return has, user, deviceInfo, err
}

// 检测激活码
func CheckActiveCode(code string) (has bool, user string, deviceInfo Device, err error) {
	has, user, deviceInfo, err = getActiveCodeValue(code)
	return has, user, deviceInfo, err
}

// 激活时，验证激活码
func VerifyActiveCode(code string) (result bool, user string, device Device) {
	has, user, deviceInfo, err := CheckActiveCode(code)
	if err == nil && has {
		status := deviceInfo.State
		if status == 1 {
			result = true
		}
	}
	return result, user, deviceInfo
}

// 重置短信状态，保证激活之后能再次收到短信，否则要等8小时之后了
func ResetSmsStatus(user, deviceId string) (err error) {
	phone := ""
	key := fmt.Sprintf("sms_%v_%v_%v", user, deviceId, phone)
	InitRedis()
	vars.RedisInstance.Del(key)
	return err
}

// 在激活流程中使用过激活码后重置激活码的状态，设为已激活并在2小时后自动删除
func ResetActiveCodeStatus(code, username string, device Device, state int) (result bool, err error) {
	InitRedis()
	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.CodePrefix, code)
	device.State = state
	err = setActiveCode(device, code)
	if err == nil {
		result, err = vars.RedisInstance.Expire(key, 60*60*2*time.Second).Result()
		err = ResetSmsStatus(username, device.DeviceId)
	}
	return result, err
}

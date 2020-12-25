package models

import (
	"encoding/json"
	"fmt"

	"exchange_zero_trust_api/util"
	"exchange_zero_trust_api/vars"
)

type (
	// 移动设备信息
	Device struct {
		DeviceType string `json:"devicetype"`
		DeviceId   string `json:"deviceid"`
		User       string `json:"user"`
		State      int    `json:"state"`
		Time       int64  `json:"allowtime"`
	}

	UserPhone struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
)

// 获取设备状态：激活、未激活、锁定等
func GetDeviceState(user, deviceId string) (state int) {
	state = -1
	InitRedis()
	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.UserPrefix, user)
	ret, err := vars.RedisInstance.HGet(key, deviceId).Result()
	if err == nil {
		var device Device
		err := json.Unmarshal([]byte(ret), &device)
		if err == nil {
			state = device.State
		}
	}
	return state
}

// 获取设备信息
func GetDeviceInfo(user, deviceId string) (err error, deviceInfo Device) {
	InitRedis()
	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.UserPrefix, user)
	info, err := vars.RedisInstance.HGet(key, deviceId).Result()
	if err == nil {
		err = json.Unmarshal([]byte(info), &deviceInfo)
	}
	return err, deviceInfo
}

// 设置设备的状态：激活，锁定
// ALLOW = 0
// NEW = 1
// LOCKED = 2
// BLOCK = 3
func SetDeviceState(user, deviceId string, state int) (err error) {
	InitRedis()
	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.UserPrefix, user)
	err, deviceInfo := GetDeviceInfo(user, deviceId)
	if err != nil {
		return err
	}

	deviceInfo.State = state
	info, err := json.Marshal(deviceInfo)
	v := make(map[string]interface{})
	v[deviceId] = string(info)

	if err == nil {
		_, err = vars.RedisInstance.HMSet(key, v).Result()
	}

	return err
}

// 激活设备
func ActiveDevice(user, deviceId string) (err error) {
	state := 0
	err = SetDeviceState(user, deviceId, state)
	return err
}

// 恢复设备
func RestoreDevice(user, deviceId string) (err error) {
	state := 1
	err = SetDeviceState(user, deviceId, state)
	return err
}

// 锁定设备
func LockDevice(user, deviceId string) (err error) {
	state := 2
	err = SetDeviceState(user, deviceId, state)
	return err
}

// 忽略设备
func IgnoreDevice(user, deviceId string) (err error) {
	state := 3
	err = SetDeviceState(user, deviceId, state)
	return err
}

// 获取设备列表
func GetDeviceList(username string) (devices []string, err error) {
	key := fmt.Sprintf("iuser_%v", username)
	InitRedis()
	devices, err = vars.RedisInstance.HVals(key).Result()
	return devices, err
}

// 获取设备数
func GetDeviceNum(username string) (n int) {
	devices, err := GetDeviceList(username)
	if err == nil {
		n = len(devices)
	}
	return n
}

// 通过设备ID查询手机设备的信息
func GetDeviceInfoByDeviceId(deviceId string) (deviceInfo util.DeviceInfo, err error) {
	InitRedis()
	key := vars.RedisKeyPrefix.DeviceInfo
	deviceInfoStr, err := vars.RedisInstance.HGet(key, deviceId).Result()
	err = json.Unmarshal([]byte(deviceInfoStr), &deviceInfo)
	return deviceInfo, err
}

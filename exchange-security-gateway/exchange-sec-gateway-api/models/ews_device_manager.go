package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"exchange_zero_trust_api/logger"
	"exchange_zero_trust_api/vars"
)

type (
	EwsCode struct {
		Username   string `json:"username"`
		Ip         string `json:"ip"`
		State      int    `json:"state"`
		ClientType string `json:"client_type"`
		IpList     string `json:"ip_list"`
	}

	EwsInfo struct {
		// Time       time.Time `json:"time"`
		// ExpireTime time.Time `json:"expire_time"`

		Time       int64  `json:"time"`
		ExpireTime int64  `json:"expire_time"`
		ClientType string `json:"client_type"`
		State      int    `json:"state"`
		EwsType    int    `json:"ews_type"`
	}

	EwsResponse struct {
		Username     string `json:"username"`
		ClientType   string `json:"client_type"`
		ClientIp     string `json:"client_ip"`
		Area         string `json:"area"`
		IpList       string `json:"ip_list"`
		ActiveStatus string `json:"active_status"`
	}

	EwsActiveResp struct {
		Username   string `json:"username"`
		ClientType string `json:"client_type"`
		ClientIp   string `json:"client_ip"`
		Area       string `json:"area"`
		IpList     string `json:"ip_list"`
		Code       string `json:"code"`
	}

	EwsTrustAddress struct {
		State   int `json:"state"`
		EwsType int `json:"ews_type"`
		// Now        time.Time
		// ExpireTime time.Time
		Now        int64  `json:"now"`
		ExpireTime int64  `json:"expire_time"`
		ClientType string `json:"client_type"`
	}
)

// 获取激活码的内容
func getEwsActiveCodeValue(code string) (*EwsCode, error) {
	var (
		ewsCode EwsCode
		err     error
	)

	key := fmt.Sprintf("EWS_CODE_%v", code)
	InitRedis()
	result, err := vars.RedisInstance.HGetAll(key).Result()
	logger.Log.Errorf("key: %v, result: %v, err: %v", key, result, err)
	if err != nil {
		return nil, err
	}

	ewsCode.Username = result["username"]
	ewsCode.Ip = result["ip"]
	state, err := strconv.Atoi(result["state"])
	ewsCode.State = state
	ewsCode.ClientType = result["client_type"]
	ewsCode.IpList = result["iplist"]

	return &ewsCode, err
}

// 判断激活码是否存在
func ExistEwsActiveCode(code string) (bool, *EwsCode, error) {
	var (
		result bool
		err    error
	)

	ewsCode, err := getEwsActiveCodeValue(code)
	if err != nil {
		return false, nil, err
	}

	if ewsCode.State >= 0 {
		result = true
	}

	return result, ewsCode, err
}

// 判断激活码是否有效
func VerifyEwsActiveCode(code string) (bool, *EwsCode, error) {
	var (
		result  bool
		err     error
		ewsCode *EwsCode
	)

	ewsCode, err = getEwsActiveCodeValue(code)
	if err != nil {
		return false, nil, err
	}

	if ewsCode.State == 0 {
		result = true
	}

	return result, ewsCode, err
}

// 判断特定用户的可信ip状态，如果state<0，表示禁用
func GetEwsIpStatus(username, ip string) (bool, error) {
	var (
		result  = true
		err     error
		ewsInfo EwsInfo
	)

	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.EwsPrefix, username)
	InitRedis()
	res, err := vars.RedisInstance.HMGet(key, ip).Result()
	logger.Log.Warnf("res: %v, err: %v", res, err)
	if err != nil {
		return false, err
	}

	if len(res) > 0 {
		info, ok := res[0].(string)
		logger.Log.Errorf("info: %v, ok: %v", info, ok)
		if ok {
			err := json.Unmarshal([]byte(info), &ewsInfo)
			if err != nil {
				return false, err
			}

			state := ewsInfo.State
			if state < 0 {
				result = false
			}
		}
	}

	return result, err
}

// 设置激活码状态，2表示已使用
func SetEwsActiveCodeState(code, username, ip string) error {
	key := fmt.Sprintf("EWS_CODE_%v", code)
	InitRedis()

	value := make(map[string]interface{})
	value["username"] = username
	value["ip"] = ip
	value["state"] = 2

	_, err := vars.RedisInstance.HMSet(key, value).Result()
	vars.RedisInstance.Expire(key, 30*time.Second)

	// reset ews sms status
	err = ResetEwsSmsStatus(username)
	return err
}

func ResetEwsSmsStatus(username string) error {
	key := fmt.Sprintf("EWS_SMS_%v", username)
	InitRedis()
	_, err := vars.RedisInstance.Del(key).Result()
	return err
}

// 添加EWS可信IP
// ews_type=0, 表示自动加入，ews_type=1，表示手动激活加入
// state=0, 表示激活，state=1, 表示未激活，state=-1, 表示禁止
func AddEwsTrustIp(username string, ip string, ewsType int, state int, clientType string) error {
	key := fmt.Sprintf("%v%v", vars.RedisKeyPrefix.EwsPrefix, username)

	InitRedis()
	now := time.Now()
	hours12, _ := time.ParseDuration("12h")
	expireTime := now.Add(hours12)

	ewsTrustAddr := EwsTrustAddress{}
	{
		ewsTrustAddr.Now = now.Unix()
		ewsTrustAddr.ExpireTime = expireTime.Unix()
		ewsTrustAddr.State = state
		ewsTrustAddr.EwsType = ewsType
		ewsTrustAddr.ClientType = clientType
	}

	addrInfo, err := json.Marshal(ewsTrustAddr)
	value := make(map[string]interface{})
	value[ip] = addrInfo

	_, err = vars.RedisInstance.HMSet(key, value).Result()

	return err
}

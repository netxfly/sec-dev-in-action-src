package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"exchange_zero_trust_api/vars"
)

type (
	TokenResp struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
	}

	MessageText struct {
		Content string `json:"content"`
	}

	MessageContent struct {
		MsgType string      `json:"msgtype"`
		Text    MessageText `json:"text"`
	}

	MessageReq struct {
		AgentId    int            `json:"agent_id"`
		UserIdList string         `json:"userid_list"`
		DeptIdList string         `json:"deptid_list"`
		ToAllUser  bool           `json:"to_all_user"`
		Msg        MessageContent `json:"msg"`
	}

	MessageResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		TaskId  int    `json:"task_id"`
	}
)

func GetAliToken() (string, error) {
	apiUrl := fmt.Sprintf("%v/gettoken?appkey=%v&appsecret=%v", vars.OapiHost, vars.AliAppKey, vars.AliAppSecret)
	client := http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequest("GET", apiUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	body, _ := ioutil.ReadAll(resp.Body)
	var respToken TokenResp
	err = json.Unmarshal(body, &respToken)
	token := respToken.AccessToken

	return token, err
}

func SendDingDingMessage(user, message string) (bool, error) {
	token, _ := GetAliToken()
	sendApi := fmt.Sprintf("%v/topapi/message/corpconversation/asyncsend_v2?access_token=%v", vars.OapiHost, token)
	client := http.Client{Timeout: 10 * time.Second}

	messageReq := MessageReq{}
	messageReq.AgentId = vars.AliAgentId
	messageReq.UserIdList = user
	// messageReq.ToAllUser = true
	messageReq.Msg.Text = MessageText{Content: message}
	messageReq.Msg.MsgType = "text"

	jsMessage, _ := json.Marshal(messageReq)
	fmt.Println(jsMessage)
	req, _ := http.NewRequest("POST", sendApi, bytes.NewBuffer(jsMessage))
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	body, _ := ioutil.ReadAll(resp.Body)
	var respMsg MessageResp
	err = json.Unmarshal(body, &respMsg)
	result := false
	if respMsg.ErrCode == 0 {
		result = true
	}
	return result, err
}

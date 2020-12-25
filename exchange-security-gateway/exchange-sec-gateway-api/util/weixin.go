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
	TokenAPI struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	SendMessageResp struct {
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
		InvalidUser  string `json:"invaliduser"`
		InvalidParty string `json:"invalidparty"`
		InvalidTag   string `json:"invalidtag"`
	}

	Text struct {
		Content string `json:"content"`
	}

	SendMessage struct {
		ToUser        string `json:"touser"`
		ToParty       string `json:"toparty"`
		ToTag         string `json:"totag"`
		MsgType       string `json:"msgtype"`
		AgentId       int    `json:"agentid"`
		Text          Text   `json:"text"`
		Safe          int    `json:"safe"`
		EnableIDTrans int    `json:"enable_id_trans"`
	}
)

func GetAccessToken() (string, error) {
	tokenApiUrl := fmt.Sprintf("%v/cgi-bin/gettoken?corpid=%v&corpsecret=%v",
		vars.CorpHost, vars.CorpId, vars.CorpSecret)

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", tokenApiUrl, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenApi TokenAPI
	err = json.Unmarshal(body, &tokenApi)
	if err != nil {
		return "", err
	}

	token := tokenApi.AccessToken
	return token, err

}

func SendWeiXin(username, content string) (bool, error) {
	result := false
	accessToken, err := GetAccessToken()
	if err != nil {
		return false, err
	}
	sendApiUrl := fmt.Sprintf("%v/cgi-bin/message/send?access_token=%v", vars.CorpHost, accessToken)
	var message SendMessage
	message.ToUser = username
	message.MsgType = "text"
	message.AgentId = vars.AgentId
	message.Text = Text{Content: content}

	client := http.Client{Timeout: 10 * time.Second}
	jsMessage, _ := json.Marshal(message)
	req, err := http.NewRequest("POST", sendApiUrl, bytes.NewBuffer(jsMessage))
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var respMessage SendMessageResp
	err = json.Unmarshal(body, &respMessage)
	if err != nil {
		return false, err
	}
	if respMessage.ErrCode == 0 {
		result = true
	}

	return result, err
}

package hw_push

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/**
 * init
 */
func NewClient(clientID string, clientSecret string) *HuaweiPushClient {
	fmt.Println("evercyan/hw_push", "NewClient", "clientID", clientID, "clientSecret", clientSecret)
	vers := &Vers{
		Ver:   "1",
		AppID: clientID,
	}
	nspCtx, _ := json.Marshal(vers)
	return &HuaweiPushClient{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		NspCtx:       string(nspCtx),
	}
}

/**
 * message init
 */
func NewMessage() *Message {
	return &Message{
		Hps: Hps{
			Msg: Msg{
				Type: 3, //1, 透传异步消息; 3, 系统通知栏异步消息;
				Body: Body{
					Content: "",
					Title:   "",
				},
				Action: Action{
					Type: 1, //1, 自定义行为; 2, 打开URL; 3, 打开App;
					Param: Param{
						Intent: "#Intent;compo=com.rvr/.Activity;S.W=U;end",
					},
				},
			},
			Ext: Ext{ // 扩展信息, 含 BI 消息统计, 特定展示风格, 消息折叠;
				Action:  "",
				Func:    "",
				Collect: "",
				Title:   "",
				Content: "",
				Url:     "",
			},
		},
	}
}

/**
 * form post
 */
func FormPost(url string, data url.Values) ([]byte, error) {
	u := ioutil.NopCloser(strings.NewReader(data.Encode()))
	r, err := http.Post(url, "application/x-www-form-urlencoded", u)
	if err != nil {
		fmt.Println("evercyan/hw_push", "FormPost", "err", err)
		return []byte(""), err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("evercyan/hw_push", "FormPost", "err", err)
		return []byte(""), err
	}
	return b, err
}

/**
 * get token
 */
func (this HuaweiPushClient) GetToken() string {
	reqUrl := TOKEN_URL
	param := make(url.Values)
	param["grant_type"] = []string{GRANTTYPE}
	param["client_id"] = []string{this.ClientId}
	param["client_secret"] = []string{this.ClientSecret}
	res, err := FormPost(reqUrl, param)
	fmt.Println("evercyan/hw_push", "GetToken", "res", string(res), "err", err)
	if nil != err {
		return ""
	}
	var tokenRes = &TokenResStruct{}
	err = json.Unmarshal(res, tokenRes)
	if err != nil {
		return ""
	}
	return tokenRes.Access_token
}

/**
 * push msg
 */
func (this HuaweiPushClient) PushMsg(deviceToken, payload string) string {
	fmt.Println("evercyan/hw_push", "PushMsg", "deviceToken", deviceToken, "payload", payload)
	accessToken := this.GetToken()
	reqUrl := PUSH_URL + "?nsp_ctx=" + url.QueryEscape(this.NspCtx)
	fmt.Println("evercyan/hw_push", "PushMsg", "reqUrl", reqUrl)

	var originParam = map[string]string{
		"access_token":      accessToken,
		"nsp_svc":           NSP_SVC,
		"nsp_ts":            strconv.Itoa(int(time.Now().Unix())),
		"device_token_list": "[\"" + deviceToken + "\"]",
		"payload":           payload,
		"expire_time":       time.Now().Format("2006-01-02T15:04"),
	}
	fmt.Println("evercyan/hw_push", "param", originParam)

	param := make(url.Values)
	param["access_token"] = []string{originParam["access_token"]}
	param["nsp_svc"] = []string{originParam["nsp_svc"]}
	param["nsp_ts"] = []string{originParam["nsp_ts"]}
	param["device_token_list"] = []string{originParam["device_token_list"]}
	param["payload"] = []string{originParam["payload"]}

	// push
	res, _ := FormPost(reqUrl, param)
	fmt.Println("evercyan/hw_push", "PushMsg", "res", string(res))
	return string(res)
}

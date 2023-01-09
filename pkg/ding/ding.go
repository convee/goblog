package ding

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	dingTalkHost = "https://oapi.dingtalk.com"
	Token        = ""
)

/**
 * 发送钉钉报警
 * token：报警机器人的token
 * content：报警内容
 * all：true：at所有人
 */

// SendAlert 发送钉钉报警
func SendAlert(content string, all bool) ([]byte, error) {
	dingUrl := dingTalkHost + "/robot/send?access_token=" + Token
	data := make(map[string]interface{})

	data["msgtype"] = "text"
	data["text"] = map[string]string{"content": "【" + time.Now().Format("2006-01-02 03:04:05") + "】" + content}
	data["at"] = map[string]interface{}{"atMobiles": [0]string{}, "isAtAll": all}

	bytePayload, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, postErr := http.Post(dingUrl, "application/json", bytes.NewBuffer(bytePayload))
	if postErr != nil {
		return nil, postErr
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

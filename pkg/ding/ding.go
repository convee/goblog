package ding

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const dingTalkHost = "https://oapi.dingtalk.com/robot/send"

var alertClient = &http.Client{Timeout: 5 * time.Second}

func SendAlert(content string, all bool) ([]byte, error) {
	token := os.Getenv("GOBLOG_DINGTALK_TOKEN")
	if token == "" {
		return nil, errors.New("GOBLOG_DINGTALK_TOKEN is not configured")
	}
	endpoint, err := url.Parse(dingTalkHost)
	if err != nil {
		return nil, err
	}
	query := endpoint.Query()
	query.Set("access_token", token)
	endpoint.RawQuery = query.Encode()

	payload, err := jsoniter.Marshal(map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "【" + time.Now().Format("2006-01-02 15:04:05") + "】" + content,
		},
		"at": map[string]interface{}{"atMobiles": []string{}, "isAtAll": all},
	})
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := alertClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return body, fmt.Errorf("dingtalk returned %s", response.Status)
	}
	return body, nil
}

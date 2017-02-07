package baidupush

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"time"
)

const (
	// DefaultBaiduPushService is the default push service host.
	DefaultBaiduPushService = "api.tuisong.baidu.com"
	// SDKNameVersion is the SDK name and version.
	SDKNameVersion = "Golang Baidu Push Service SDK v1.0"
	// MsgTypeMessage represents a push of message.
	MsgTypeMessage = 0
	// MsgTypeNotice represents a push of notification.
	MsgTypeNotice = 1
	// AndroidDeviceType represents Android platform number for Baidu Push Service.
	AndroidDeviceType = 3
	// AppleDeviceType represents Apple platform number for Baidu Push Service.
	AppleDeviceType = 4
	// DeployStatusDevelop represents development status.
	DeployStatusDevelop = 1
	// DeployStatusProduct represents production status.
	DeployStatusProduct = 2
)

// Channel contains all the methods to interact with Baidu Cloud Push Service.
type Channel struct {
	host       string
	apiKey     string
	secret     string
	requestID  int64
	deviceType int
}

// NewChannel returns a channel bound with specified paramters.
func NewChannel(host, key, secret string, device int) *Channel {
	return &Channel{
		host:       host,
		apiKey:     key,
		secret:     secret,
		deviceType: device,
	}
}

// NewChannelDefaultHost returns a channel with host set to "api.tuisong.baidu.com"
func NewChannelDefaultHost(key, secret string, device int) *Channel {
	return NewChannel(DefaultBaiduPushService, key, secret, device)
}

// GetRequestID returns request ID returned by server.
func (bc *Channel) GetRequestID() int64 {
	return bc.requestID
}

func absorbOptionalKeys(commons, opts url.Values) url.Values {
	together := url.Values{}
	for k, v := range commons {
		together[k] = v
	}
	for k, v := range opts {
		together[k] = v
	}
	return together
}

func requestService(host, apiClass, apiMethod, httpMethod, secret string, query url.Values) ([]byte, error) {
	urlStr := fmt.Sprintf("http://%s/rest/3.0/%s/%s", host, apiClass, apiMethod)
	sign := generateSign(httpMethod, urlStr, secret, query)
	query.Add("sign", sign)

	req, err := http.NewRequest(httpMethod, urlStr, bytes.NewReader([]byte(query.Encode())))
	if err != nil {
		return nil, err
	}

	req.Header = apiHeader()
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	return ioutil.ReadAll(rsp.Body)
}

// PushMsgToSingleDevice pushes a message to a single device.
func (bc *Channel) PushMsgToSingleDevice(channelID string, msg string, opts url.Values) (string, int64, error) {
	var msgID string
	var sendTime int64

	err := checkOptionalKeys("PushMsgToSingleDevice", opts)
	if err != nil {
		return msgID, sendTime, err
	}

	query := absorbOptionalKeys(commonRequestParams(bc.apiKey, bc.deviceType), opts)
	query.Add("channel_id", channelID)
	query.Add("msg", msg)

	data, err := requestService(bc.host, "push", "single_device", http.MethodPost, bc.secret, query)
	if err != nil {
		return msgID, sendTime, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return msgID, sendTime, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return msgID, sendTime, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})
	msgID = rspParams["msg_id"].(string)
	sendTime = int64(rspParams["send_time"].(float64))
	return msgID, sendTime, nil
}

// PushMsgToAllDevice pushes a message to all devices running app.
func (bc *Channel) PushMsgToAllDevice(msg string, opts url.Values) (string, string, int64, error) {
	var msgID, timerID string
	var sendTime int64

	err := checkOptionalKeys("PushMsgToAllDevice", opts)
	if err != nil {
		return msgID, timerID, sendTime, err
	}
	query := absorbOptionalKeys(commonRequestParams(bc.apiKey, bc.deviceType), opts)
	query.Add("msg", msg)

	data, err := requestService(bc.host, "push", "all", http.MethodPost, bc.secret, query)
	if err != nil {
		return msgID, timerID, sendTime, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return msgID, timerID, sendTime, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return msgID, timerID, sendTime, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})
	msgID = rspParams["msg_id"].(string)
	sendTime = int64(rspParams["send_time"].(float64))
	if tid, ok := rspParams["timer_id"]; ok {
		timerID = tid.(string)
	}
	return msgID, timerID, sendTime, nil
}

func commonRequestParams(apiKey string, deviceType int) url.Values {
	commons := url.Values{}
	commons.Add("apikey", apiKey)
	commons.Add("timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	commons.Add("device_type", fmt.Sprintf("%d", deviceType))
	return commons
}

func generateSign(method, urlStr, secret string, params url.Values) string {
	gather := ""
	gather += method
	gather += urlStr

	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		gather += fmt.Sprintf("%s=%s", k, params.Get(k))
	}

	gather += secret
	data := md5.Sum([]byte(url.QueryEscape(gather)))
	return hex.EncodeToString(data[:])
}

func apiHeader() http.Header {
	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	userAgent := fmt.Sprintf("BCCS_SDK/3.0 (%s) %s (%s) cli/Unknown", runtime.GOOS, runtime.Version(), SDKNameVersion)
	header.Add("User-Agent", userAgent)
	return header
}

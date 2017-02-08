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
	"strconv"
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

func (bc *Channel) pushMessage(apiName, apiMethod string, musts, optionals url.Values) (map[string]interface{}, error) {
	err := checkOptionalKeys(apiName, optionals)
	if err != nil {
		return nil, err
	}

	query := absorbOptionalKeys(commonRequestParams(bc.apiKey, bc.deviceType), musts, optionals)

	data, err := requestService(bc.host, "push", apiMethod, http.MethodPost, bc.secret, query)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return nil, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})
	resultMap := map[string]interface{}{
		"msg_id":   rspParams["msg_id"].(string),
		"timer_id": "",
	}

	switch sendTime := rspParams["send_time"].(type) {
	case string:
		val, _ := strconv.ParseInt(sendTime, 10, 64)
		resultMap["send_time"] = val
	case float64:
		val := int64(sendTime)
		resultMap["send_time"] = val
	}

	if tid, ok := rspParams["timer_id"]; ok {
		resultMap["timer_id"] = tid.(string)
	}
	return resultMap, nil
}

// PushMsgToSingleDevice pushes a message to a single device.
func (bc *Channel) PushMsgToSingleDevice(channelID string, msg string, opts url.Values) (string, int64, error) {
	var msgID string
	var sendTime int64

	musts := url.Values{}
	musts.Add("channel_id", channelID)
	musts.Add("msg", msg)

	resultMap, err := bc.pushMessage("PushMsgToSingleDevice", "single_device", musts, opts)
	if err != nil {
		return msgID, sendTime, err
	}

	msgID = resultMap["msg_id"].(string)
	sendTime = resultMap["send_time"].(int64)

	return msgID, sendTime, nil
}

// PushMsgToBatchDevices pushes a message to a batch of devices.
func (bc *Channel) PushMsgToBatchDevices(channelIDs []string, msg string, opts url.Values) (string, int64, error) {
	var msgID string
	var sendTime int64

	channelsData, err := json.Marshal(channelIDs)
	if err != nil {
		return msgID, sendTime, err
	}

	musts := url.Values{}
	musts.Add("channel_ids", string(channelsData))
	musts.Add("msg", msg)

	resultMap, err := bc.pushMessage("PushMsgToBatchDevices", "batch_device", musts, opts)
	if err != nil {
		return msgID, sendTime, err
	}

	msgID = resultMap["msg_id"].(string)
	sendTime = resultMap["send_time"].(int64)

	return msgID, sendTime, nil
}

// PushMsgToAllDevices pushes a message to all devices running app.
func (bc *Channel) PushMsgToAllDevices(msg string, opts url.Values) (string, string, int64, error) {
	var msgID, timerID string
	var sendTime int64

	musts := url.Values{}
	musts.Add("msg", msg)

	resultMap, err := bc.pushMessage("PushMsgToAllDevice", "all", musts, opts)
	if err != nil {
		return msgID, timerID, sendTime, err
	}

	msgID = resultMap["msg_id"].(string)
	timerID = resultMap["timer_id"].(string)
	sendTime = resultMap["send_time"].(int64)

	return msgID, timerID, sendTime, nil
}

// PushMsgToTaggedDevices pushes a message to devices under some tag.
func (bc *Channel) PushMsgToTaggedDevices(tag, msg string, opts url.Values) (string, string, int64, error) {
	var msgID, timerID string
	var sendTime int64

	musts := url.Values{}
	musts.Add("type", fmt.Sprintf("%d", 1))
	musts.Add("tag", tag)
	musts.Add("msg", msg)

	resultMap, err := bc.pushMessage("PushMsgToTag", "tags", musts, opts)
	if err != nil {
		return msgID, timerID, sendTime, err
	}

	msgID = resultMap["msg_id"].(string)
	timerID = resultMap["timer_id"].(string)
	sendTime = resultMap["send_time"].(int64)

	return msgID, timerID, sendTime, nil
}

// MessageResult represents the information about sent message.
type MessageResult struct {
	MsgID    string
	Status   int
	Success  int
	SendTime int64
}

// QueryMsgStatus queries message reports via msgID.
func (bc *Channel) QueryMsgStatus(msgID string) (int, []MessageResult, error) {
	totalNum := 0

	musts := url.Values{}
	musts.Add("msg_id", msgID)

	resultMap, err := bc.query("QueryMsgStatus", "query_msg_status", musts, nil)
	if err != nil {
		return totalNum, nil, err
	}

	totalNum = resultMap["total_num"].(int)
	results := resultMap["result"].([]MessageResult)

	return totalNum, results, nil
}

// QueryTimerRecords queries records of timed message via timerID.
func (bc *Channel) QueryTimerRecords(timerID string, opts url.Values) (string, []MessageResult, error) {
	var retTimerID string

	musts := url.Values{}
	musts.Add("timer_id", timerID)

	resultMap, err := bc.query("QueryTimerRecords", "query_timer_records", musts, opts)
	if err != nil {
		return retTimerID, nil, err
	}

	retTimerID = resultMap["timer_id"].(string)
	results := resultMap["result"].([]MessageResult)
	return retTimerID, results, nil
}

// QueryTopicRecords queries records of topic message via topicID.
func (bc *Channel) QueryTopicRecords(topicID string, opts url.Values) (string, []MessageResult, error) {
	var retTopicID string

	musts := url.Values{}
	musts.Add("topic_id", topicID)

	resultMap, err := bc.query("QueryTopicRecords", "query_topic_records", musts, opts)
	if err != nil {
		return retTopicID, nil, err
	}

	retTopicID = resultMap["topic_id"].(string)
	results := resultMap["result"].([]MessageResult)
	return retTopicID, results, nil
}

// TimerResult represents information about timed task.
type TimerResult struct {
	ID        string
	Msg       string
	SendTime  int64
	MsgType   int
	RangeType int
}

// QueryTimerTasks queries timer tasks not executing yet.
func (bc *Channel) QueryTimerTasks(opts url.Values) (int, []TimerResult, error) {
	var totalNum int

	err := checkOptionalKeys("QueryTimerTasks", opts)
	if err != nil {
		return totalNum, nil, err
	}

	query := absorbOptionalKeys(commonRequestParams(bc.apiKey, bc.deviceType), opts)

	data, err := requestService(bc.host, "timer", "query_list", http.MethodGet, bc.secret, query)
	if err != nil {
		return totalNum, nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return totalNum, nil, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return totalNum, nil, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})
	totalNum = int(rspParams["total_num"].(float64))

	timerResults := []TimerResult{}
	resultData := rspParams["result"].([]interface{})
	for _, r := range resultData {
		rMap := r.(map[string]interface{})
		timerResult := TimerResult{
			ID:        rMap["timer_id"].(string),
			Msg:       rMap["msg"].(string),
			SendTime:  int64(rMap["send_time"].(float64)),
			MsgType:   int(rMap["msg_type"].(float64)),
			RangeType: int(rMap["range_type"].(float64)),
		}
		timerResults = append(timerResults, timerResult)
	}

	return totalNum, timerResults, nil
}

// CancelTimerTask cancels timed message not executing yet.
func (bc *Channel) CancelTimerTask(timerID string) error {
	query := commonRequestParams(bc.apiKey, bc.deviceType)
	query.Add("timer_id", timerID)

	data, err := requestService(bc.host, "timer", "cancel", http.MethodPost, bc.secret, query)
	if err != nil {
		return err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return checkErrorCode(int(errCode.(float64)))
	}

	return nil
}

func (bc *Channel) query(apiName, apiMethod string, musts, optionals url.Values) (map[string]interface{}, error) {
	err := checkOptionalKeys(apiName, optionals)
	if err != nil {
		return nil, err
	}

	query := absorbOptionalKeys(commonRequestParams(bc.apiKey, bc.deviceType), musts, optionals)

	data, err := requestService(bc.host, "report", apiMethod, http.MethodGet, bc.secret, query)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return nil, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})

	resultMap := map[string]interface{}{}
	if total, ok := rspParams["total_num"]; ok {
		resultMap["total_num"] = int(total.(float64))
	}

	if timer, ok := rspParams["timer_id"]; ok {
		resultMap["timer_id"] = timer.(string)
	}

	if topic, ok := rspParams["topic_id"]; ok {
		resultMap["topic_id"] = topic.(string)
	}

	results := []MessageResult{}
	resultData := rspParams["result"].([]interface{})
	for _, r := range resultData {
		rMap := r.(map[string]interface{})
		queryResult := MessageResult{
			MsgID:    rMap["msg_id"].(string),
			Status:   int(rMap["status"].(float64)),
			SendTime: int64(rMap["send_time"].(float64)),
		}
		if succ, ok := rMap["success"]; ok {
			queryResult.Success = int(succ.(float64))
		}
		results = append(results, queryResult)
	}

	resultMap["result"] = results

	return resultMap, nil
}

// TagInfo represents information about a tag.
type TagInfo struct {
	TID        string
	Tag        string
	Info       string
	Type       int // deprecated
	CreateTime int64
}

// QueryTagsInfo querys tags information of app, opts contains optional parameters below.
//
// tag: name of the tag, must be of length 1-128, "default" is reserved so cannot be used.
//
// start: the start position of the returned records, defaults to 0.
//
// limit: the number of records returned, must be 1-100, defaults to 100.
func (bc *Channel) QueryTagsInfo(opts url.Values) (int, []TagInfo, error) {
	totalNum := 0
	tagInfos := []TagInfo{}

	err := checkOptionalKeys("QueryTagsInfo", opts)
	if err != nil {
		return totalNum, nil, err
	}
	query := absorbOptionalKeys(commonRequestParams(bc.apiKey, bc.deviceType), opts)

	data, err := requestService(bc.host, "app", "query_tags", http.MethodGet, bc.secret, query)
	if err != nil {
		return totalNum, nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return totalNum, nil, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return totalNum, nil, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})
	totalNum = int(rspParams["total_num"].(float64))
	tags := rspParams["result"].([]interface{})

	for _, tagData := range tags {
		tagMap := tagData.(map[string]interface{})
		tid := tagMap["tid"].(string)
		tag := tagMap["tag"].(string)
		info := tagMap["info"].(string)
		typ := int(tagMap["type"].(float64))
		createTime := int64(tagMap["create_time"].(float64))
		tagInfo := TagInfo{
			TID:        tid,
			Tag:        tag,
			Info:       info,
			Type:       typ,
			CreateTime: createTime,
		}
		tagInfos = append(tagInfos, tagInfo)
	}
	return totalNum, tagInfos, nil
}

// CreateTag creates an empty tag group.
//
// tag: name of the tag, must be of length 1-128, "default" is reserved so cannot be used.
func (bc *Channel) CreateTag(tag string) (string, error) {
	return bc.manageTag("create_tag", tag)
}

// DeleteTag deletes an existed tag group.
//
// tag: name of the tag, must be of length 1-128, "default" is reserved so cannot be used.
func (bc *Channel) DeleteTag(tag string) (string, error) {
	return bc.manageTag("del_tag", tag)
}

func (bc *Channel) manageTag(apiMethod, tag string) (string, error) {
	retTag := ""

	query := commonRequestParams(bc.apiKey, bc.deviceType)
	query.Add("tag", tag)

	data, err := requestService(bc.host, "app", apiMethod, http.MethodPost, bc.secret, query)
	if err != nil {
		return retTag, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return retTag, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return retTag, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})
	retTag = rspParams["tag"].(string)
	retCode := int(rspParams["result"].(float64))
	if retCode != 0 {
		return retTag, fmt.Errorf("code %d - create tag failed", retCode)
	}

	return retTag, nil
}

// TagResult represents the result of add/delete devices from tag group.
type TagResult struct {
	ChnID string
	Res   int
}

// AddTagDevices adds a batch of devices to a tag group.
//
// tag: name of the tag, must be of length 1-128, "default" is reserved so cannot be used.
//
// channelIDs: a string slice containing channel IDs to add, require at least 1 and at most 10.
func (bc *Channel) AddTagDevices(tag string, channelIDs []string) ([]TagResult, error) {
	return bc.manageTagDevices("add_devices", tag, channelIDs)
}

// DeleteTagDevices deletes a batch of devices from a tag group.
//
// tag: name of the tag, must be of length 1-128, "default" is reserved so cannot be used.
//
// channelIDs: a string slice containing channel IDs to add, require at least 1 and at most 10.
func (bc *Channel) DeleteTagDevices(tag string, channelIDs []string) ([]TagResult, error) {
	return bc.manageTagDevices("del_devices", tag, channelIDs)
}

// GetTagDevicesNumber returns the number of devices related to tag.
func (bc *Channel) GetTagDevicesNumber(tag string) (int, error) {
	num := 0

	query := commonRequestParams(bc.apiKey, bc.deviceType)
	query.Add("tag", tag)

	data, err := requestService(bc.host, "tag", "device_num", http.MethodGet, bc.secret, query)
	if err != nil {
		return num, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return num, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return num, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})
	num = int(rspParams["device_num"].(float64))

	return num, nil
}

func (bc *Channel) manageTagDevices(apiMethod, tag string, channelIDs []string) ([]TagResult, error) {
	tagResults := []TagResult{}

	if len(channelIDs) < 1 || len(channelIDs) > 10 {
		return nil, fmt.Errorf("invalid channel ID number %d - must be [1, 10]", len(channelIDs))
	}

	chnData, err := json.Marshal(channelIDs)
	if err != nil {
		return nil, err
	}

	query := commonRequestParams(bc.apiKey, bc.deviceType)
	query.Add("tag", tag)
	query.Add("channel_ids", string(chnData))

	data, err := requestService(bc.host, "tag", apiMethod, http.MethodPost, bc.secret, query)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	bc.requestID = int64(result["request_id"].(float64))
	if errCode, ok := result["error_code"]; ok {
		return nil, checkErrorCode(int(errCode.(float64)))
	}

	rspParams := result["response_params"].(map[string]interface{})
	devices := rspParams["result"].([]interface{})
	for _, dev := range devices {
		devMap := dev.(map[string]interface{})
		cid := devMap["channel_id"].(string)
		res := int(devMap["result"].(float64))
		tagResults = append(tagResults, TagResult{ChnID: cid, Res: res})
	}

	return tagResults, nil
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

func absorbOptionalKeys(valuesSet ...url.Values) url.Values {
	together := url.Values{}

	for _, values := range valuesSet {
		for k, v := range values {
			together[k] = v
		}
	}

	return together
}

func requestService(host, apiClass, apiMethod, httpMethod, secret string, query url.Values) ([]byte, error) {
	urlStr := fmt.Sprintf("http://%s/rest/3.0/%s/%s", host, apiClass, apiMethod)
	sign := generateSign(httpMethod, urlStr, secret, query)
	query.Add("sign", sign)

	var req *http.Request
	var err error
	if httpMethod == http.MethodPost {
		req, err = http.NewRequest(httpMethod, urlStr, bytes.NewReader([]byte(query.Encode())))
		if err != nil {
			return nil, err
		}
	} else if httpMethod == http.MethodGet {
		urlStr = fmt.Sprintf("%s?%s", urlStr, query.Encode())
		req, err = http.NewRequest(httpMethod, urlStr, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header = apiHeader()
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	return ioutil.ReadAll(rsp.Body)
}

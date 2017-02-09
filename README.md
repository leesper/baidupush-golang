BaiduPush-Golang Library v1.0
===========================================

A Golang SDK for working with Baidu Cloud Push Service.

Package baidupush provides a third-party Golang SDK for working with Baidu Cloud Push Service.

Channel contains all the functions for working with the service, including pushing, tagging, querying and timed tasks.

For detailed information, please referring the official documentation: [Baidu](http://push.baidu.com/document)

# Documentation
```go
import "github.com/leesper/couchdb-golang"
```

## Constants
```go
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
```

## type Channel
```go
type Channel struct {
    // contains filtered or unexported fields
}
```
Channel contains all the methods to interact with Baidu Cloud Push Service.

## func NewChannel
```go
func NewChannel(host, key, secret string, device int) *Channel
```
NewChannel returns a channel bound with specified paramters.

host: URL address of Baidu Cloud Push Service.

key: API key.

secret: API secret.

device: Device type, AppleDeviceType or AndroidDeviceType.

## func NewChannelDefaultHost
```go
func NewChannelDefaultHost(key, secret string, device int) *Channel
```
NewChannelDefaultHost returns a channel with host set to "api.tuisong.baidu.com"

## func (\*Channel) AddTagDevices
```go
func (bc *Channel) AddTagDevices(tag string, channelIDs []string) ([]TagResult, error)
```
AddTagDevices adds a batch of devices to a tag group.

tag: name of the tag, must be of length 1-128, "default" is reserved so cannot be used.

channelIDs: a string slice containing channel IDs to add, require at least 1 and at most 10.

## func (\*Channel) CancelTimerTask
```go
func (bc *Channel) CancelTimerTask(timerID string) error
```
CancelTimerTask cancels timed message not executing yet.

timerID: ID of timed task.

## func (\*Channel) CreateTag
```
func (bc *Channel) CreateTag(tag string) (string, error)
```
CreateTag creates an empty tag group.

tag: Name of the tag, must be of length 1-128, "default" is reserved so cannot be used.

## func (\*Channel) DeleteTag
```go
func (bc *Channel) DeleteTag(tag string) (string, error)
```
DeleteTag deletes an existed tag group.

tag: Name of the tag, must be of length 1-128, "default" is reserved so cannot be used.

## func (\*Channel) DeleteTagDevices
```go
func (bc *Channel) DeleteTagDevices(tag string, channelIDs []string) ([]TagResult, error)
```
DeleteTagDevices deletes a batch of devices from a tag group.

tag: name of the tag, must be of length 1-128, "default" is reserved so cannot be used.

channelIDs: a string slice containing channel IDs to add, require at least 1 and at most 10.

## func (\*Channel) GetRequestID
```go
func (bc *Channel) GetRequestID() int64
```
GetRequestID returns request ID returned by server.

## func (\*Channel) GetTagDevicesNumber
```go
func (bc *Channel) GetTagDevicesNumber(tag string) (int, error)
```
GetTagDevicesNumber returns the number of devices related to tag.

## func (\*Channel) PushMsgToAllDevices
```go
func (bc *Channel) PushMsgToAllDevices(msg string, opts url.Values) (string, string, int64, error)
```
PushMsgToAllDevices pushes a message to all devices running app.

msg: Message to push

Optional parameters:

msg_type: Type of message, MsgTypeNotice or MsgTypeMessage(default).

msg_expires: 0-604800, defaults to 5 hrs, the time message expired(from the moment on).

deploy_status: Deployment status(for iOS app only), DeployStatusProduct(default) or DeployStatusDevelop.

send_time: The real sending time for timed message, must be at least 60s and at most 1 year.

## func (\*Channel) PushMsgToBatchDevices
```go
func (bc *Channel) PushMsgToBatchDevices(channelIDs []string, msg string, opts url.Values) (string, int64, error)
```
PushMsgToBatchDevices pushes a message to a batch of devices.

channelIDs: Channel IDs of devices.

msg: Message to push.

Optional parameters:

msg_type: Type of message, MsgTypeNotice or MsgTypeMessage(default).

msg_expires: 0-604800, defaults to 5 hrs, the time message expired(from the moment on).

topic_id: Name of the topic.

## func (\*Channel) PushMsgToSingleDevice
```go
func (bc *Channel) PushMsgToSingleDevice(channelID string, msg string, opts url.Values) (string, int64, error)
```
PushMsgToSingleDevice pushes a message to a single device.

channelID: Channel ID of the device.

msg: Message to push.

Optional parameters:

msg_type: Type of message, MsgTypeNotice or MsgTypeMessage(default).

msg_expires: 0-604800, defaults to 5 hrs, the time message expired(from the moment on).

deploy_status: Deployment status(for iOS app only), DeployStatusProduct(default) or DeployStatusDevelop.

## func (\*Channel) PushMsgToTaggedDevices
```go
func (bc *Channel) PushMsgToTaggedDevices(tag, msg string, opts url.Values) (string, string, int64, error)
```
PushMsgToTaggedDevices pushes a message to devices under some tag.

tag: Name of created tag.

msg: Message to push.

Optional parameters:

msg_type: Type of message, MsgTypeNotice or MsgTypeMessage(default).

msg_expires: 0-604800, defaults to 5 hrs, the time message expired(from the moment on).

deploy_status: Deployment status(for iOS app only), DeployStatusProduct(default) or DeployStatusDevelop.

send_time: The real sending time for timed message, must be at least 60s and at most 1 year.

## func (\*Channel) QueryMsgStatus
```go
func (bc *Channel) QueryMsgStatus(msgID string) (int, []MessageResult, error)
```
QueryMsgStatus queries message reports via msgID.

msgID: Message ID, could be a json array of IDs.

## func (\*Channel) QueryTagsInfo
```go
func (bc *Channel) QueryTagsInfo(opts url.Values) (int, []TagInfo, error)
```
QueryTagsInfo querys tags information of app, opts contains optional parameters below.

Optional parameters:

tag: name of the tag, must be of length 1-128, "default" is reserved so cannot be used.

start: the start position of the returned records, defaults to 0.

limit: the number of records returned, must be 1-100, defaults to 100.

## func (\*Channel) QueryTimerRecords
```go
func (bc *Channel) QueryTimerRecords(timerID string, opts url.Values) (string, []MessageResult, error)
```
QueryTimerRecords queries records of timed message via timerID.

timerID: ID of timer task.

Optional parameters:

start: the start position of the returned records, defaults to 0.

limit: the number of records returned, must be 1-100, defaults to 100.

range_start: UNIX timestamp, the start time to query.

range_end: UNIX timestamp, the end time to query.

## func (\*Channel) QueryTimerTasks
```go
func (bc *Channel) QueryTimerTasks(opts url.Values) (int, []TimerResult, error)
```
QueryTimerTasks queries timer tasks not executing yet.

Optional parameters:

timer_id: ID of timed task.

start: the start position of the returned records, defaults to 0.

limit: the number of records returned, must be 1-100, defaults to 100.

## func (\*Channel) QueryTopicList
```go
func (bc *Channel) QueryTopicList(opts url.Values) (int, []TopicResult, error)
```
QueryTopicList returns topics been used.

Optional parameters:

start: the start position of the returned records, defaults to 0.

limit: the number of records returned, must be 1-100, defaults to 100.

## func (\*Channel) QueryTopicRecords
```go
func (bc *Channel) QueryTopicRecords(topicID string, opts url.Values) (string, []MessageResult, error)
```
QueryTopicRecords queries records of topic message via topicID.

topicID: Name of the topic.

Optional parameters:

start: the start position of the returned records, defaults to 0.

limit: the number of records returned, must be 1-100, defaults to 100.

range_start: UNIX timestamp, the start time to query.

range_end: UNIX timestamp, the end time to query.

## func (\*Channel) ReportDeviceStatistics
```go
func (bc *Channel) ReportDeviceStatistics() (int, []DeviceStatistics, error)
```
ReportDeviceStatistics returns statistics about devices installed app.

## func (\*Channel) ReportTopicStatistics
```go
func (bc *Channel) ReportTopicStatistics(topicID string) (int, []TopicStatistics, error)
```
ReportTopicStatistics returns statistic information about number of messages under some topic.

topicID: Name of the topic.

## type DeviceStatistics
```go
type DeviceStatistics struct {
    Day           int64
    DailyNewUser  int
    DailyLostUser int
    DailyOnline   int
    AddedupTerm   int
    AvailChnID    int
}
```
DeviceStatistics represents statistic about devices installed app.

## type MessageResult
```go
type MessageResult struct {
    MsgID    string
    Status   int
    Success  int
    SendTime int64
}
```
MessageResult represents the information about sent message.

## type TagInfo
```go
type TagInfo struct {
    TID        string
    Tag        string
    Info       string
    Type       int // deprecated
    CreateTime int64
}
```
TagInfo represents information about a tag.

## type TagResult
```go
type TagResult struct {
    ChnID string
    Res   int
}
```
TagResult represents the result of add/delete devices from tag group.

## type TimerResult
```go
type TimerResult struct {
    ID        string
    Msg       string
    SendTime  int64
    MsgType   int
    RangeType int
}
```
TimerResult represents information about timed task.

## type TopicResult
```go
type TopicResult struct {
    AckCount, PushCount int
    FirstTime, LastTime int64
    Topic               string
}
```
TopicResult represents information of topic.

## type TopicStatistics
```go
type TopicStatistics struct {
    Day int64
    Ack int
}
```
TopicStatistics represents statistic information about topic.

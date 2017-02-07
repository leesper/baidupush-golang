package baidupush_test

import (
	"fmt"
	"net/url"

	baidupush "github.com/leesper/baidupush-golang"
)

var (
	channelID = "4215667327923129295"
	apiKey    = "NM2AmKF7f84qw7l26h1ICEVf"
	secret    = "LlITGDusuKhLTuBKoBjP2yCZLql3Ieun"
	channel   = baidupush.NewChannelDefaultHost(apiKey, secret, baidupush.AndroidDeviceType)
)

func ExampleChannel_pushMsgToSingleDevice() {
	msg := `{
    "title": "hello",
    "description": "hello world",
    "notification_basic_style": 7,
    "open_type": 1,
    "url": "http://developer.baidu.com",
    "custom_content":{"push":"single_device"}
	}`

	opts := url.Values{}
	msgType := fmt.Sprintf("%d", baidupush.MsgTypeNotice)
	opts.Add("msg_type", msgType)
	/*msgID, sendTime, err := */ channel.PushMsgToSingleDevice(channelID, msg, opts)
	// fmt.Println(msgID, sendTime, err)
	fmt.Println("msgID", "sendTime", "err")
	// Output:
	// msgID sendTime err
}

func ExampleChannel_pushMsgToAllDevice() {
	msg := `{
    "title": "hello",
    "description": "hello world",
    "notification_basic_style": 7,
    "open_type": 1,
    "url": "http://developer.baidu.com",
    "custom_content":{"push":"all"}
	}`

	opts := url.Values{}
	msgType := fmt.Sprintf("%d", baidupush.MsgTypeNotice)
	opts.Add("msg_type", msgType)
	/*msgID, timerID, sendTime, err := */ channel.PushMsgToAllDevice(msg, opts)
	// fmt.Println(msgID, timerID, sendTime, err)
	fmt.Println("msgID", "timerID", "sendTime", "err")
	// Output:
	// msgID timerID sendTime err
}

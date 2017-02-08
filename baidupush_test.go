package baidupush

import (
	"fmt"
	"net/url"
	"testing"
)

var (
	channelID = "4215667327923129295"
	apiKey    = "NM2AmKF7f84qw7l26h1ICEVf"
	secret    = "LlITGDusuKhLTuBKoBjP2yCZLql3Ieun"
	channel   = NewChannelDefaultHost(apiKey, secret, AndroidDeviceType)
)

func TestPushMsgToSingleDevice(t *testing.T) {
	msg := `{
    "title": "hello",
    "description": "hello world",
    "notification_basic_style": 7,
    "open_type": 1,
    "url": "http://developer.baidu.com",
    "custom_content":{"push":"single_device"}
	}`

	opts := url.Values{}
	msgType := fmt.Sprintf("%d", MsgTypeNotice)
	opts.Add("msg_type", msgType)
	msgID, sendTime, err := channel.PushMsgToSingleDevice(channelID, msg, opts)
	if err != nil {
		t.Fatal("push to single error", err)
	}

	fmt.Printf("push to single success, request ID %d message ID %s time %d\n",
		channel.GetRequestID(), msgID, sendTime)
}

func TestPushMsgToAllDevice(t *testing.T) {
	msg := `{
    "title": "hello",
    "description": "hello world",
    "notification_basic_style": 7,
    "open_type": 1,
    "url": "http://developer.baidu.com",
    "custom_content":{"push":"all"}
	}`

	opts := url.Values{}
	msgType := fmt.Sprintf("%d", MsgTypeNotice)
	opts.Add("msg_type", msgType)
	msgID, timerID, sendTime, err := channel.PushMsgToAllDevices(msg, opts)
	if err != nil {
		t.Fatal("push to all error", err)
	}

	fmt.Printf("push to all success, request ID %d message ID %s timer ID %s time %d\n",
		channel.GetRequestID(), msgID, timerID, sendTime)
}

func TestTagManagement(t *testing.T) {
	tag1, err := channel.CreateTag("tag1")
	if err != nil {
		t.Error("create tag1 error", err)
	}
	if tag1 != "tag1" {
		t.Errorf("create tag returns %s want tag1", tag1)
	}

	tag2, err := channel.CreateTag("tag2")
	if err != nil {
		t.Error("create tag2 error", err)
	}
	if tag2 != "tag2" {
		t.Errorf("create tag returns %s want tag2", tag2)
	}

	tag3, err := channel.CreateTag("tag3")
	if err != nil {
		t.Error("create tag3 error", err)
	}
	if tag3 != "tag3" {
		t.Errorf("create tag returns %s want tag3", tag3)
	}

	totalNum, tagsInfo, err := channel.QueryTagsInfo(nil)
	if err != nil {
		t.Error("query tags info error", err)
	}

	if totalNum != 4 { // tag1 tag2 tag3 and default
		t.Errorf("total number %d want 4", totalNum)
	}

	for _, info := range tagsInfo {
		if info.Tag != "default" && info.Tag != "tag1" && info.Tag != "tag2" && info.Tag != "tag3" {
			t.Errorf("tag %s want tag1 or tag2 or tag3", info.Tag)
		}
	}

	results, err := channel.AddTagDevices("tag2", []string{channelID})
	if err != nil {
		t.Error("add tag2 devices error", err)
	} else {
		if len(results) != 1 {
			t.Errorf("len(results) = %d want 1", len(results))
		}
		if results[0].ChnID != channelID {
			t.Errorf("return channel ID %s want %s", results[0].ChnID, channelID)
		}
		if results[0].Res != 0 {
			t.Error("add devices to tag failed")
		}
	}

	_, err = channel.GetTagDevicesNumber("tag2")
	if err != nil {
		t.Error("get tag devices number error", err)
	}

	results, err = channel.DeleteTagDevices("tag2", []string{channelID})
	if err != nil {
		t.Error("add tag2 devices error", err)
	} else {
		if len(results) != 1 {
			t.Errorf("len(results) = %d want 1", len(results))
		}
		if results[0].ChnID != channelID {
			t.Errorf("return channel ID %s want %s", results[0].ChnID, channelID)
		}
	}

	tag1, err = channel.DeleteTag("tag1")
	if err != nil {
		t.Error("delete tag1 error", err)
	}
	if tag1 != "tag1" {
		t.Errorf("delete tag returns %s want tag1", tag1)
	}

	tag2, err = channel.DeleteTag("tag2")
	if err != nil {
		t.Error("delete tag2 error", err)
	}
	if tag2 != "tag2" {
		t.Errorf("delete tag returns %s want tag2", tag2)
	}

	tag3, err = channel.DeleteTag("tag3")
	if err != nil {
		t.Error("delete tag3 error", err)
	}
	if tag3 != "tag3" {
		t.Errorf("delete tag returns %s want tag3", tag3)
	}

}

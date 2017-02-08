package baidupush

import (
	"fmt"
	"net/url"
)

var (
	optionalKeys = map[string]map[string]bool{
		"PushMsgToSingleDevice": {
			"expires":       true,
			"device_type":   true,
			"msg_type":      true,
			"msg_expires":   true,
			"deploy_status": true,
		},
		"PushMsgToAllDevice": {
			"expires":       true,
			"device_type":   true,
			"msg_type":      true,
			"msg_expires":   true,
			"deploy_status": true,
			"send_time":     true,
		},
		"QueryTagsInfo": {
			"expires":     true,
			"device_type": true,
			"tag":         true,
			"start":       true,
			"limit":       true,
		},
		"PushMsgToTag": {
			"expires":       true,
			"device_type":   true,
			"msg_type":      true,
			"msg_expires":   true,
			"deploy_status": true,
			"send_time":     true,
		},
		"PushMsgToBatchDevices": {
			"expires":     true,
			"device_type": true,
			"msg_type":    true,
			"msg_expires": true,
			"topic_id":    true,
		},
	}
)

func checkOptionalKeys(api string, params url.Values) error {
	optionals, ok := optionalKeys[api]
	if !ok {
		return fmt.Errorf("API %s not found", api)
	}
	for k := range params {
		if _, ok = optionals[k]; !ok {
			return fmt.Errorf("invalid parameter: %s is not allowed in API %s", k, api)
		}
	}
	return nil
}

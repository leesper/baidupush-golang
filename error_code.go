package baidupush

import "fmt"

var (
	errCodeMap = map[int]error{
		30600: fmt.Errorf("%d - %s", 30600, "internal server error"),
		30601: fmt.Errorf("%d - %s", 30601, "method not allowed"),
		30602: fmt.Errorf("%d - %s", 30602, "request params not valid"),
		30603: fmt.Errorf("%d - %s", 30603, "authentication failed"),
		30604: fmt.Errorf("%d - %s", 30604, "quota use up, payment required"),
		30605: fmt.Errorf("%d - %s", 30605, "data required not found"),
		30606: fmt.Errorf("%d - %s", 30606, "request time expires timeout"),
		30607: fmt.Errorf("%d - %s", 30607, "channel token timeout"),
		30608: fmt.Errorf("%d - %s", 30608, "bind relation not found"),
		30609: fmt.Errorf("%d - %s", 30609, "bind number too many"),
		30610: fmt.Errorf("%d - %s", 30610, "duplicate operation"),
		30611: fmt.Errorf("%d - %s", 30611, "tag not found"),
		30612: fmt.Errorf("%d - %s", 30612, "app forbidden, need whitelist authorization"),
		30613: fmt.Errorf("%d - %s", 30613, "app need initiated first in push console"),
		30616: fmt.Errorf("%d - %s", 30616, "app is not approved, can not use the push service"),
		30617: fmt.Errorf("%d - %s", 30617, "app do not have broadcast push capability"),
		30618: fmt.Errorf("%d - %s", 30618, "app do not have unicast or groupcast push capability"),
		30619: fmt.Errorf("%d - %s", 30619, "default tag is reserved"),
		30620: fmt.Errorf("%d - %s", 30620, "one app could only have one kind of device platform"),
		30621: fmt.Errorf("%d - %s", 30621, "package name invalid"),
		30699: fmt.Errorf("%d - %s", 30699, "requests are too frequent to be temporarily rejected or need whitelist authorization"),
		40001: fmt.Errorf("%d - %s", 40001, "invalid iOS device token"),
		40002: fmt.Errorf("%d - %s", 40002, "invalid iOS message"),
		40003: fmt.Errorf("%d - %s", 40003, "iOS bad device token"),
		40004: fmt.Errorf("%d - %s", 40004, "iOS certification error"),
		40005: fmt.Errorf("%d - %s", 40005, "iOS duplicate message"),
		40006: fmt.Errorf("%d - %s", 40006, "iOS production certification invalid"),
		40007: fmt.Errorf("%d - %s", 40007, "iOS development certification invalid"),
		40008: fmt.Errorf("%d - %s", 40008, "iOS production certification expire"),
		40009: fmt.Errorf("%d - %s", 40009, "iOS development certification expire"),
		40010: fmt.Errorf("%d - %s", 40010, "type error, need a development certification"),
		40011: fmt.Errorf("%d - %s", 40011, "type error, need a production certification"),
		40012: fmt.Errorf("%d - %s", 40012, "iOS certification file invalid"),
		41001: fmt.Errorf("%d - %s", 41001, "timer task not exist"),
		41002: fmt.Errorf("%d - %s", 41002, "timer task duplicated"),
		41003: fmt.Errorf("%d - %s", 41003, "timer task num exceed"),
		41004: fmt.Errorf("%d - %s", 41004, "timer task will be executed, can not be canceled"),
		41005: fmt.Errorf("%d - %s", 41005, "timer task has been executed"),
		50001: fmt.Errorf("%d - %s", 50001, "generate CSRF token failed"),
		50002: fmt.Errorf("%d - %s", 50002, "invalid CSRF token"),
		50003: fmt.Errorf("%d - %s", 50003, "CSRF token expired"),
		50004: fmt.Errorf("%d - %s", 50004, "passport not login"),
		50005: fmt.Errorf("%d - %s", 50005, "invalid BDUSS"),
		50006: fmt.Errorf("%d - %s", 50006, "required to register as a developer"),
		50007: fmt.Errorf("%d - %s", 50007, "invalid developer"),
		50008: fmt.Errorf("%d - %s", 50008, "invalid app name"),
	}
)

func checkErrorCode(code int) error {
	if err, ok := errCodeMap[code]; ok {
		return err
	}
	return nil
}

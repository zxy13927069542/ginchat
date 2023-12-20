package utils

import (
	"github.com/asaskevich/govalidator"
	"regexp"
)

func Init() {
	govalidator.TagMap["mobile"] = IsMobile
}

//	IsMobile 校验是否是手机号码
func IsMobile(s string) bool {
	matched, _ := regexp.MatchString(`^(1[1-9][0-9]\d{8})$`, s)
	return matched
}

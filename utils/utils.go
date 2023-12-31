package utils

import (
	"github.com/asaskevich/govalidator"
	"regexp"
	"strconv"
)

func Init() {
	govalidator.TagMap["mobile"] = IsMobile
}

//	IsMobile 校验是否是手机号码
func IsMobile(s string) bool {
	matched, _ := regexp.MatchString(`^(1[1-9][0-9]\d{8})$`, s)
	return matched
}

func Str2Uint(s string) uint {
	v, _ := strconv.Atoi(s)
	return uint(v)
}

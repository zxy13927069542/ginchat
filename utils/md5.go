package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Md5Encode() md5加密后转成16进制
func Md5Encode(raw string) string {
	h := md5.New()
	h.Write([]byte(raw))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5Encode() 加密后转成大写
func MD5Encode(raw string) string {
	return strings.ToUpper(Md5Encode(raw))
}

// Sum 生成md5校验和
func Sum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

// MakePassword() 生成md5加密加盐的密码
func MakePassword(plain string, salt string) string {
	return Md5Encode(plain + salt)
}

// ValidatePassword() 校验密码
func ValidatePassword(plain string, salt string, password string) bool {
	return MakePassword(plain, salt) == password
}

func GenSalt() string {
	return fmt.Sprintf("%03v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000))
}

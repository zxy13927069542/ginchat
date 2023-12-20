package utils

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	tokenStr := GenToken("郑小燕", "066311")
	token, _ := ParseToken(tokenStr.Token)
	fmt.Printf("token.UserName = %s\n", token.UserName)
	fmt.Printf("token.Password = %s\n", token.Password)
	if token.UserName != "郑小燕" || token.Password != "066311" {
		t.Fail()
	}
}

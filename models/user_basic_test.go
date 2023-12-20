package models

import (
	"fmt"
	"ginchat/config"
	"testing"
)

func init() {
	c := config.Init("../etc")
	_ = Init(c)
}

func TestList(t *testing.T) {
	list, _ := NewUserBasicModel().List()
	for _, v := range list {
		fmt.Println(v)
	}
}

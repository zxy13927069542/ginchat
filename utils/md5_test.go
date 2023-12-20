package utils

import (
	"fmt"
	"testing"
)

func TestGenSalt(t *testing.T) {
	fmt.Printf("生成的随机盐: %v\n", GenSalt())
}

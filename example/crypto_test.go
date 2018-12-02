package example_test

import (
	"fmt"
	"testing"

	"git.code.oa.com/fip-team/fiorm"
)

// TestHashPassword 密码测试
func TestHashPassword(t *testing.T) {
	pwd := "pass01"
	hashKey := fiorm.HashPassword(pwd)
	fmt.Println(hashKey)

	ok := fiorm.CheckPassword(hashKey, pwd)
	fmt.Println(ok)
}

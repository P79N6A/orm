package example_test

import (
	"fmt"
	"reflect"
	"testing"

	"git.code.oa.com/fip-team/fiorm"
	"git.code.oa.com/fip-team/fiorm/model"
)

// 创建表

func TestCreateTable(t *testing.T) {
	var user model.User
	var dept model.Department

	fiorm.CreateTable(&user)
	fiorm.CreateTable(&dept)

}

func TestAa(t *testing.T) {
	var user int32
	msg := fmt.Sprintf("未知的ID类型,%v", reflect.TypeOf(user))
	fmt.Println(msg)
}

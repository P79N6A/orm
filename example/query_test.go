package example_test

import (
	"fmt"
	"testing"

	"git.code.oa.com/fip-team/fiorm"
	"git.code.oa.com/fip-team/fiorm/model"
)

func TestGetItemByID(t *testing.T) {
	var user model.User

	da := fiorm.DataAccess()
	da.GetItemByID(&user, 46)

	fmt.Println(user.Name)
}

func TestGetItemWhereFirst(t *testing.T) {
	var user model.User

	da := fiorm.DataAccess()
	query := fiorm.Where("name =?", "wins")

	da.GetItemWhereFirst(&user, query)

	fmt.Println(user.Name)
}

func TestGetItemWhere(t *testing.T) {
	var users []model.User

	da := fiorm.DataAccess()
	query := fiorm.Where("name =?", "wins").
		OrderBy("id desc").
		Limit(3)

	da.GetItemWhere(&users, query)

	for _, user := range users {
		println(user.ID)
	}
}

func TestConvertToDto(t *testing.T) {
	// UserDto 为自定义的类，包含多个表字段
	var results []model.UserDto

	da := fiorm.DataAccess()
	da.ExecuteTextQuery(&results,
		"SELECT u.id,u.name, dept_name FROM department t ,user u WHERE u.dept_id=t.id AND u.name =? and u.address like ?",
		"wins", "inse%")

	for _, userDto := range results {
		fmt.Println(userDto.DeptName)
	}
}

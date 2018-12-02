package example_test

import (
	"fmt"
	"testing"

	"git.code.oa.com/fip-team/fiorm"
	"git.code.oa.com/fip-team/fiorm/model"
)

func TestScanToDto(t *testing.T) {
	var results []model.UserDto

	da := fiorm.DataAccess()
	da.ExecuteTextQuery(&results,
		"SELECT u.id,u.name, dept_name FROM department t ,user u WHERE u.dept_id=t.id AND u.name =? and u.address like ?",
		"wins", "inse%")

	for _, userDto := range results {
		fmt.Println(userDto.DeptName)
	}
}

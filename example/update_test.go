package example_test

import (
	"strconv"
	"testing"

	"github.com/guregu/null"

	"git.code.oa.com/fip-team/fiorm"

	"git.code.oa.com/fip-team/fiorm/model"
)

func TestUpdateOne(t *testing.T) {
	var user model.User
	da := fiorm.DataAccess()
	da.GetItemByID(&user, 160940)
	user.Email = null.StringFrom("121212121@qq.com")
	user.Address = null.StringFrom("12121")

	da.UpdateItem(&user, []string{"email", "address"})
}

func TestBatchUpdate(t *testing.T) {
	var users []model.User
	da := fiorm.DataAccess()
	query := fiorm.Where("name=?", "wins")
	da.GetItemWhere(&users, query)

	for idx, _ := range users {
		users[idx].Address = null.StringFrom("aaa" + strconv.Itoa(idx))
		users[idx].Email = null.StringFrom("22@11.com" + strconv.Itoa(idx))
	}

	// 批量更新，表定义必须包含一个ID的主键
	da.UpdateItem(users, []string{"email", "address"})
}

package example_test

import (
	"testing"

	"git.code.oa.com/fip-team/fiorm"

	"git.code.oa.com/fip-team/fiorm/model"
)

func TestDeleteOne(t *testing.T) {
	var user model.User
	da := fiorm.DataAccess()
	da.GetItemByID(&user, 6)

	da.DeleteItem(&user)
}

func TestDeleteMany(t *testing.T) {
	var users []model.User

	da := fiorm.DataAccess()
	query := fiorm.Where("name =?", "wins").
		OrderBy("id desc").
		Limit(3)

	da.GetItemWhere(&users, query)

	da.DeleteItem(&users)
}

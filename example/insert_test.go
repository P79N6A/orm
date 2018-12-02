package example_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"git.code.oa.com/fip-team/fiorm"
	"github.com/guregu/null"

	"git.code.oa.com/fip-team/fiorm/model"
)

// 单行创建
func TestInsertOne(t *testing.T) {
	var user model.User

	user.Address = null.StringFrom("深圳")
	user.CreatedAt = null.TimeFrom(time.Now())
	user.DeptID = null.IntFrom(2)
	//user.Email = ""
	user.Name = null.StringFrom("wins")
	//user.NullAge = null.IntFrom(21)
	user.NullString = null.StringFrom("")
	user.Birthday = null.TimeFrom(time.Now())
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	da := fiorm.DataAccess()
	da.InsertItem(&user)

	var dept model.Department
	dept.DeptName = null.StringFrom("企业IT")
	dept.Tel = "100011123"

	da.InsertItem(&dept)

}
func TestInsertRows(t *testing.T) {
	var users []model.User
	for i := 0; i < 10000; i++ {
		var user model.User

		user.Address = null.StringFrom("深圳-" + strconv.Itoa(i))
		user.CreatedAt = null.TimeFrom(time.Now())
		user.DeptID = null.IntFrom(2)
		//user.Email = ""
		user.Name = null.StringFrom("wins")
		user.NullAge = null.IntFrom(21)
		//user.NullString = ""
		user.Birthday = null.TimeFrom(time.Now())

		users = append(users, user)
	}

	da := fiorm.DataAccess()
	da.InsertItem(users)
}

// 事务处理
func TestInsertWithTranction(t *testing.T) {
	var users []model.User
	for i := 0; i < 10000; i++ {
		var user model.User

		user.Address = null.StringFrom("深圳-" + strconv.Itoa(i))
		user.CreatedAt = null.TimeFrom(time.Now())
		user.DeptID = null.IntFrom(2)
		//user.Email = ""
		user.Name = null.StringFrom("wins")
		user.NullAge = null.IntFrom(21)
		//user.NullString = ""
		user.Birthday = null.TimeFrom(time.Now())

		users = append(users, user)
	}

	start := time.Now().String()

	//da:=fiorm.DataAccess()
	//da.InsertItem(users)

	tx := fiorm.BeginTranction()

	tx.InsertItem(users)
	tx.InsertItem(users)
	tx.EndTranction()
	// 2万行耗时2.8s
	fmt.Println(start, "---", time.Now())
}

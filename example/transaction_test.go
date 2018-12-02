package example_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"git.code.oa.com/fip-team/fiorm"
	"git.code.oa.com/fip-team/fiorm/model"
	"github.com/guregu/null"
)

//事务测试
func TestSimpleTx(t *testing.T) {
	var user model.User
	user.Address = null.StringFrom("深圳")
	user.CreatedAt = null.TimeFrom(time.Now())
	user.DeptID = null.IntFrom(2)
	//user.Email = ""
	user.Name = null.StringFrom("wins")
	//user.NullAge = null.IntFrom(21)
	user.NullString = null.StringFrom("")
	user.Birthday = null.TimeFrom(time.Now())

	var user2 model.User
	user2.Address = null.StringFrom("深圳")
	user2.CreatedAt = null.TimeFrom(time.Now())
	user2.DeptID = null.IntFrom(2)
	//user2.Email = ""
	user2.Name = null.StringFrom("wins")
	//user2.NullAge = null.IntFrom(21)
	user2.NullString = null.StringFrom("")
	user2.Birthday = null.TimeFrom(time.Now())

	var user3 model.User
	user3.Address = null.StringFrom("深圳")
	user3.CreatedAt = null.TimeFrom(time.Now())
	user3.DeptID = null.IntFrom(2)
	//user3.Email = ""
	user3.Name = null.StringFrom("wins")
	//user3.NullAge = null.IntFrom(21)
	user3.NullString = null.StringFrom("")
	user3.Birthday = null.TimeFrom(time.Now())

	// 当事务错误时，可以在这里处理错误信息
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error handle here")
		}
	}()

	fmt.Println("begin tranction")
	tx := fiorm.BeginTranction()

	tx.InsertItem(&user3)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go insert(tx, &user, &wg, 0)
	}

	wg.Add(1)
	go insert(tx, &user2, &wg, 0)

	wg.Wait()
	tx.EndTranction()
	fmt.Println("end tranction")

}

func insert(tx *fiorm.FiDB, value interface{}, wg *sync.WaitGroup, idx int) {
	tx.InsertItem(value)
	wg.Done()

	if idx == 1 {
		panic("测试错误")
	}
}

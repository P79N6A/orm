package fiorm

import (
	"fmt"
	"reflect"
	"strconv"

	"git.code.oa.com/fip-team/fiorm/internal"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DB 返回一个连接池的实例
var db *gorm.DB

// FiDB fiorm处理类
type FiDB struct {
	db          *gorm.DB // 局部
	Error       error
	isTranction bool // 最外层是否已经开启事务，如果开启，批量插入或更新则取消
}

// DbSetting 数据库连接字符串属性
type DbSettings struct {
	Dialect  string
	DbName   string
	Host     string
	User     string
	Password string
	Port     int
}

// InitDB 初始化
//
// dialect -别名,如mysql,oracle
// InitDB("mysql", "GO_TESTDB", "localhost", "root", "password", 3306)
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
func InitDB(conf *DbSettings) {
	pwd, err1 := internal.Decrypt(conf.Password)
	if err1 != nil {
		panic("密码解密错误")
	}

	internal.DbSetting.Dialect = conf.Dialect
	internal.DbSetting.DbName = conf.DbName
	internal.DbSetting.Host = conf.Host
	internal.DbSetting.User = conf.User
	internal.DbSetting.Password = conf.Password
	internal.DbSetting.Port = conf.Port

	p := strconv.Itoa(conf.Port)
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=true&loc=Local", conf.User, pwd, conf.Host, p, conf.DbName)
	if source == "" {
		panic("错误的连接字符串")
	}

	var err error
	db, err = gorm.Open(conf.Dialect, source)
	if err != nil {
		panic("数据库没有初始化")
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	db.SingularTable(true)
	db.LogMode(true)

	// 设置日志
	lg := new(GoLogWriter)
	db.SetLogger(gorm.Logger{lg})

}

// DataAccess 获取一个数据库连接
func DataAccess() *FiDB {
	var da = new(FiDB)

	if db == nil {
		panic("DB不能为空")
	}

	da.db = db.New()
	return da
}

// GetGormDB 获取Gorm原生DB
//func GetGormDB() *gorm.DB {
//	return db
//}

// GetItemByID 根据主键ID获取数据
func (t *FiDB) GetItemByID(tEntity interface{}, id int64) {
	t.db = t.db.First(tEntity, id)

	if t.db.Error != nil && t.db.Error != gorm.ErrRecordNotFound {
		panic(t.db.Error)
	}
}

// GetItemWhereFirst 根据条件查询一条数据
func (t *FiDB) GetItemWhereFirst(tEntity interface{}, query *Query) {
	if query.db.Error != nil {
		panic(query.db.Error)
	}

	t.db = query.db.First(tEntity)
}

// GetItemWhere 根据条件查询多条数据
func (t *FiDB) GetItemWhere(tEntity interface{}, query *Query) {
	if query.db.Error != nil {
		panic(query.db.Error)
	}

	t.db = query.db.Find(tEntity)
}

// Count 返回总行数
func (t *FiDB) Count(value interface{}) {
	if t.db.Error != nil && t.db.Error != gorm.ErrRecordNotFound {
		panic(t.db.Error)
	}

	t.db.Count(value)
}

// InsertItem 插入一条或批量
func (t *FiDB) InsertItem(value interface{}) {
	kind := reflect.TypeOf(value).Kind()
	if kind == reflect.Ptr || kind == reflect.Struct {
		t.db = t.db.Create(value)
	} else {
		// 启动事务，防止部分写入
		if t.isTranction == false {
			tx := BeginTranction()
			tx.batchInsert(value)
			tx.EndTranction()
		} else {
			// 如果外部已经启动事务，这里需要取消事务，否则将导致部分更新成功
			t.batchInsert(value)
		}
	}

	if t.db.Error != nil {
		panic(t.db.Error)
	}
}

// DeleteItem 删除一条或多条数据
func (t *FiDB) DeleteItem(value interface{}) {
	// TODO 判断要删除的实体是否存在，如果不存在，需要返回错误信息
	if t.db.Error == gorm.ErrRecordNotFound {
		return
	}
	// 防止删除所有数据
	if t.db.Error != nil {
		panic(t.db.Error)
	}

	// 必须用t.db而不是全局db,全局db删除的情况下，如果value没有值，将导致删除所有的数据！
	t.db.Delete(value)
}

// UpdateItem 更新一行或多行数据
func (t *FiDB) UpdateItem(value interface{}, cols []string) {
	if t.db.Error == gorm.ErrRecordNotFound {
		return
	}
	// 防止更新所有数据
	if t.db.Error != nil {
		panic(t.db.Error)
	}

	kind := reflect.TypeOf(value).Kind()
	if kind == reflect.Ptr || kind == reflect.Struct {
		// 必须用t.db而不是全局db,全局db更新的情况下，如果value没有值，将导致更新所有的数据！
		t.db.Select(cols).Model(value).Updates(value)
	} else {
		// 启动事务，防止部分更新
		if t.isTranction == false {
			tx := BeginTranction()
			tx.batchUpdate(value, cols)
			tx.EndTranction()
		} else {
			// 如果外部已经启动事务，这里需要取消事务，否则将导致部分更新成功
			t.batchUpdate(value, cols)
		}
	}
}

// CreateTable 创建表 DDL操作
func CreateTable(value interface{}) {
	db.CreateTable(value)
}

// BeginTranction 开始事务
func BeginTranction() *FiDB {
	var tx = new(FiDB)
	tx.db = db.Begin()
	tx.isTranction = true
	return tx
}

// EndTranction 结束事务
func (t *FiDB) EndTranction() {
	if t.db.Error != nil {
		t.db.Rollback()
		panic(t.db.Error)
	}

	t.db.Commit()
}

// ExecuteTextQuery 原生语法查询
func (t *FiDB) ExecuteTextQuery(dto interface{}, sql string, values ...interface{}) {
	d := t.db.Raw(sql, values...)
	if d.Error != nil || t.db.Error != nil {
		panic(d.Error)
	}

	d.Scan(dto)
}

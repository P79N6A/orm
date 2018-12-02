package fiorm

import "github.com/jinzhu/gorm"

// Query 查询条件
type Query struct {
	db *gorm.DB //局部
}

// Where 查询条件
func Where(query interface{}, args ...interface{}) *Query {
	qu := &Query{}
	qu.db = db.Where(query, args...)

	return qu
}

// OrderBy 排序
func (q *Query) OrderBy(value interface{}, reorder ...bool) *Query {
	gormDB := q.db.Order(value, reorder...)
	q.db = gormDB

	return q
}

// Select 选择字段
func (q *Query) Select(query interface{}, args ...interface{}) *Query {
	gormDB := q.db.Select(query, args...)
	q.db = gormDB

	return q
}

// Limit 限制行数
func (q *Query) Limit(limit interface{}) *Query {
	gormDB := q.db.Limit(limit)
	q.db = gormDB

	return q
}

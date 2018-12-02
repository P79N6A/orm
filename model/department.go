package model

import "github.com/guregu/null"

// Department 部门表
type Department struct {
	Address   null.String `gorm:"column:address"`
	CreatedAt null.Time   `gorm:"column:created_at"`
	DeptName  null.String `gorm:"column:dept_name"`
	ID        int64       `gorm:"column:id"`
	NullFloat null.Float  `gorm:"column:null_float"`
	Tel       string      `gorm:"column:tel"`
	UpdatedAt null.Time   `gorm:"column:updated_at"`
}

// TableName sets the insert table name for this struct type
func (d *Department) TableName() string {
	return "department"
}

package model

import "github.com/guregu/null"

// User 用户表
type User struct {
	Address    null.String `gorm:"column:address"`
	Birthday   null.Time   `gorm:"column:birthday"`
	CreatedAt  null.Time   `gorm:"column:created_at"`
	DeptID     null.Int    `gorm:"column:dept_id"`
	Email      null.String `gorm:"column:email"`
	ID         int64       `gorm:"column:id"`
	Name       null.String `gorm:"column:name"`
	NullAge    null.Int    `gorm:"column:null_age"`
	NullString null.String `gorm:"column:null_string"`
	UpdatedAt  null.Time   `gorm:"column:updated_at"`
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "user"
}

package model

// UserInfoView 视图
type UserInfoView struct {
	ID       uint64
	Name     string
	DeptName string
}

func (UserInfoView) TableName() string {
	return "user_info_v"
}

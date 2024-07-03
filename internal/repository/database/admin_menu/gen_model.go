package admin_menu

import "time"

// AdminMenu 管理员菜单栏表
//
//go:generate gormgen -structs AdminMenu -input .
type AdminMenu struct {
	Id          int       // 主键
	AdminId     int       // 管理员ID
	MenuId      int       // 菜单栏ID
	CreatedAt   time.Time `gorm:"time"` // 创建时间
	CreatedUser string    // 创建人
}

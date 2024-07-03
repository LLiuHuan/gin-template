// Package admin
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 15:02
package admin

type SearchMyActionData struct {
	AdminId int `json:"admin_id"` // 管理员ID
}

type MyActionData struct {
	Id     int    // 主键
	MenuId int    // 菜单栏ID
	Method string // 请求方式
	Api    string // 请求地址
}

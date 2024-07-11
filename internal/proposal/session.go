// Package proposal
//
//	@program:		gin-template
//	@author:		[lliuhuan](https://github.com/lliuhuan)
//	@create:		2024-07-02 21:41
//	@description:	用户会话信息
package proposal

import "encoding/json"

// SessionUserInfo 当前用户会话信息
type SessionUserInfo struct {
	UserID   int    `json:"user_id"`   // 用户ID
	UserName string `json:"user_name"` // 用户名
}

// Marshal 序列化到JSON
func (user *SessionUserInfo) Marshal() (jsonRaw []byte) {
	jsonRaw, _ = json.Marshal(user)
	return
}

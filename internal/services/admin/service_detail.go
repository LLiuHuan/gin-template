// Package admin
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-09 10:03
package admin

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/database/admin"
)

type SearchOneData struct {
	Id       int    // 用户ID
	Username string // 用户名
	Nickname string // 昵称
	Mobile   string // 手机号
	Password string // 密码
	IsUsed   int32  // 是否启用 1:是  -1:否
}

func (s *service) Detail(ctx core.Context, searchOneData *SearchOneData) (info *admin.Admin, err error) {
	qb := admin.NewQueryBuilder()
	qb.WhereIsDeleted(database.EqualPredicate, -1)

	if searchOneData.Id != 0 {
		qb.WhereId(database.EqualPredicate, searchOneData.Id)
	}

	if searchOneData.Username != "" {
		qb.WhereUsername(database.EqualPredicate, searchOneData.Username)
	}

	if searchOneData.Nickname != "" {
		qb.WhereNickname(database.EqualPredicate, searchOneData.Nickname)
	}

	if searchOneData.Mobile != "" {
		qb.WhereMobile(database.EqualPredicate, searchOneData.Mobile)
	}

	if searchOneData.Password != "" {
		qb.WherePassword(database.EqualPredicate, searchOneData.Password)
	}

	if searchOneData.IsUsed != 0 {
		qb.WhereIsUsed(database.EqualPredicate, searchOneData.IsUsed)
	}

	info, err = qb.QueryOne(s.db.GetDB().WithContext(ctx.RequestContext()))
	if err != nil {
		key := ctx.GetCtx().ClientIP()
		s.cache.Incr(key)
		return nil, err
	}

	return
}

// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 11:31
package cron

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/database/cron_task"
)

func (s *service) PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error) {
	qb := cron_task.NewQueryBuilder()

	if searchData.Name != "" {
		qb.WhereName(database.EqualPredicate, searchData.Name)
	}

	if searchData.Protocol != 0 {
		qb.WhereProtocol(database.EqualPredicate, searchData.Protocol)
	}

	if searchData.IsUsed != 0 {
		qb.WhereIsUsed(database.EqualPredicate, searchData.IsUsed)
	}

	total, err = qb.Count(s.db.GetDB().WithContext(ctx.RequestContext()))
	if err != nil {
		return 0, err
	}

	return
}

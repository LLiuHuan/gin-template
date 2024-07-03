// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 11:26
package cron

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/database/cron_task"
)

type SearchOneData struct {
	Id int // 任务ID
}

func (s *service) Detail(ctx core.Context, searchOneData *SearchOneData) (info *cron_task.CronTask, err error) {
	qb := cron_task.NewQueryBuilder()

	if searchOneData.Id != 0 {
		qb.WhereId(database.EqualPredicate, searchOneData.Id)
	}

	info, err = qb.QueryOne(s.db.GetDB().WithContext(ctx.RequestContext()))
	if err != nil {
		return nil, err
	}

	return
}

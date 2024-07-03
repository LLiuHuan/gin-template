// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 11:28
package cron

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/database/cron_task"
)

func (s *service) Execute(ctx core.Context, id int) (err error) {
	qb := cron_task.NewQueryBuilder()
	qb.WhereId(database.EqualPredicate, id)
	info, err := qb.QueryOne(s.db.GetDB().WithContext(ctx.RequestContext()))
	if err != nil {
		return err
	}

	info.Spec = "手动执行"
	go s.cronServer.AddJob(info)()

	return nil
}

// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 14:15
package cron

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/database/cron_task"

	"github.com/robfig/cron/v3"
)

// TODO: 感觉cron这块整个不太对，需要重新整理

func (s *service) UpdateUsed(ctx core.Context, id int, used int32) (err error) {
	data := map[string]interface{}{
		"is_used":      used,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := cron_task.NewQueryBuilder()
	qb.WhereId(database.EqualPredicate, id)
	err = qb.Updates(s.db.GetDB().WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	// region 操作定时任务 避免主从同步延迟，在这需要查询主库
	if used == cron_task.IsUsedNo {
		s.cronServer.RemoveTask(cron.EntryID(id))
	} else {
		qb = cron_task.NewQueryBuilder()
		qb.WhereId(database.EqualPredicate, id)
		info, err := qb.QueryOne(s.db.GetDB().WithContext(ctx.RequestContext()))
		if err != nil {
			return err
		}

		s.cronServer.RemoveTask(cron.EntryID(id))
		s.cronServer.AddTask(info)

	}
	// endregion

	return
}

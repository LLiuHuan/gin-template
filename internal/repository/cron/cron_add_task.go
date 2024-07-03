// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 10:31
package cron

import (
	"strings"

	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/database/cron_task"
	"github.com/LLiuHuan/gin-template/pkg/errors"
)

func (s *server) AddTask(task *cron_task.CronTask) error {
	spec := "0 " + strings.TrimSpace(task.Spec)

	enterId, err := s.cron.AddFunc(spec, s.AddJob(task))
	if err != nil {
		return errors.New(err.Error())
	}

	//task.TaskId = enterId

	qb := cron_task.NewQueryBuilder()
	qb.WhereId(database.EqualPredicate, task.Id)
	if err := qb.Updates(s.db.GetDB(), map[string]interface{}{
		"task_id": enterId,
	}); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

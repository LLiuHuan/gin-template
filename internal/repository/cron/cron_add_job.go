// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 10:20
package cron

import (
	"fmt"

	"github.com/LLiuHuan/gin-template/internal/repository/gorm/cron_task"

	"github.com/robfig/cron/v3"
)

func (s *server) AddJob(task *cron_task.CronTask) cron.FuncJob {
	return func() {
		s.taskCount.Add()
		defer s.taskCount.Done()

		// 将 task 信息写入到 Kafka Topic 中，任务执行器订阅 Topic 如果为符合条件的任务并进行执行，反之不执行
		// 为了便于演示，不写入到 Kafka 中，仅记录日志

		msg := fmt.Sprintf("执行任务：(%d)%s [%s]", task.Id, task.Name, task.Spec)
		s.logger.Info(msg)
	}
}

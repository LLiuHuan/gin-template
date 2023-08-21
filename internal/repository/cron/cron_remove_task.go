// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-21 12:12
package cron

import (
	"github.com/robfig/cron/v3"
)

func (s *server) RemoveTask(taskId cron.EntryID) {
	s.cron.Remove(taskId)
}

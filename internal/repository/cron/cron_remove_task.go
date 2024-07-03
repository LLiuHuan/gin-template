// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 10:34
package cron

import "github.com/robfig/cron/v3"

func (s *server) RemoveTask(taskId cron.EntryID) {
	s.cron.Remove(taskId)
}

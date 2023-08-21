// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-21 12:11
package cron

func (s *server) Stop() {
	s.cron.Stop()
	s.taskCount.Exit()
}

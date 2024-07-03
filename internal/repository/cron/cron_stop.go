// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 10:35
package cron

func (s *server) Stop() {
	s.cron.Stop()
	s.taskCount.Exit()
}

// Package socket
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-04 17:35
package socket

import "go.uber.org/zap"

func (s *server) OnClose() {
	err := s.socket.Close()
	if err != nil {
		s.logger.Error("socket on closed error", zap.Error(err))
	}
}

// Package socket
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-08 14:16
package socket

func (s *server) WriteMessage(messageType int, data []byte) error {
	return s.socket.WriteMessage(messageType, data)
}

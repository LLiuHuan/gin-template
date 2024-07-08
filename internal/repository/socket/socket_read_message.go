// Package socket
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-08 14:15
package socket

func (s *server) ReadMessage() (messageType int, p []byte, err error) {
	return s.socket.ReadMessage()
}

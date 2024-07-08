// Package socket
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-08 14:18
package socket

func (s *server) WriteJSON(v interface{}) error {
	return s.socket.WriteJSON(v)
}

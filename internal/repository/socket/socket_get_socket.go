// Package socket
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-08 10:36
package socket

import "github.com/gorilla/websocket"

func (s *server) GetSocket() *websocket.Conn {
	return s.socket
}

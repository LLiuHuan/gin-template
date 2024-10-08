// Package socket
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-04 17:34
package socket

import (
	"net/http"
	"time"

	"github.com/LLiuHuan/gin-template/internal/repository/gorm"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/pkg/errors"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var _ Server = (*server)(nil)

type server struct {
	logger *zap.Logger
	db     gorm.Repo
	cache  redis.Repo
	socket *websocket.Conn
}

type Server interface {
	i()

	// OnMessage 接收消息
	OnMessage()

	// OnSend 发送消息
	OnSend(message []byte) error

	// OnClose 关闭
	OnClose()

	// GetSocket 获取socket
	GetSocket() *websocket.Conn

	ReadMessage() (messageType int, p []byte, err error)

	WriteMessage(messageType int, data []byte) error

	WriteJSON(v interface{}) error
}

var upGrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New(logger *zap.Logger, db gorm.Repo, cache redis.Repo, w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	if db == nil {
		return nil, errors.New("db required")
	}

	if cache == nil {
		return nil, errors.New("cache required")
	}

	ws, err := upGrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, errors.Wrap(err, "ws error")
	}

	return &server{
		logger: logger,
		db:     db,
		cache:  cache,
		socket: ws,
	}, nil
}

func (s *server) i() {}

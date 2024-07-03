// Package core
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 14:36
// @description: 多路复用器
package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var _ Mux = (*mux)(nil)

// Mux http mux
type Mux interface {
	http.Handler
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type mux struct {
	engine *gin.Engine
}

func (m *mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.engine.ServeHTTP(w, req)
}

func (m *mux) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &router{
		group: m.engine.Group(relativePath, wrapHandlers(handlers...)...),
	}
}

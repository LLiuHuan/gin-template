// Package generator
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-22 03:15
package generator

import "github.com/LLiuHuan/gin-template/internal/pkg/core"

func (h *handler) HandlerView() core.HandlerFunc {
	return func(c core.Context) {
		c.HTML("generator_handler", nil)
	}
}

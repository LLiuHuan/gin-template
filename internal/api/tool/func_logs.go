// Package tool
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-23 16:48
package tool

import (
	"encoding/json"
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/pkg/file"
	"go.uber.org/zap"
)

type logsParseData struct {
	Level        string  `json:"level"`
	Time         string  `json:"time"`
	Caller       string  `json:"caller"`
	Msg          string  `json:"msg"`
	Domain       string  `json:"domain"`
	Method       string  `json:"method"`
	Path         string  `json:"path"`
	HTTPCode     int     `json:"http_code"`
	BusinessCode int     `json:"business_code"`
	Success      bool    `json:"success"`
	CostSeconds  float64 `json:"cost_seconds"`
	TraceID      string  `json:"trace_id"`
}

type logsData struct {
	Level       string  `json:"level"`
	Time        string  `json:"time"`
	Path        string  `json:"path"`
	HTTPCode    int     `json:"http_code"`
	Method      string  `json:"method"`
	Msg         string  `json:"msg"`
	TraceID     string  `json:"trace_id"`
	Content     string  `json:"content"`
	CostSeconds float64 `json:"cost_seconds"`
}

type logsResponse struct {
	Logs []logsData `json:"logs"`
}

// Logs 获取日志
//
//	@Summary		获取日志
//	@Description	获取日志
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		logsRequest	true	"请求信息"
//	@Success		200		{object}	logsResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/logs [get]
func (h *handler) Logs() core.HandlerFunc {
	return func(ctx core.Context) {
		res := new(logsResponse)
		readLineFromEnd, err := file.NewReadLineFromEnd(configs.ProjectLogFile)
		if err != nil {
			h.logger.Error("NewReadLineFromEnd err", zap.Error(err))
		}

		logSize := 100
		res.Logs = make([]logsData, 100)

		index := 0

		for {
			if line, readErr := readLineFromEnd.ReadLine(); readErr == nil {
				if string(line) != "" {
					var logsParse logsParseData
					err = json.Unmarshal(line, &logsParse)
					if err != nil {
						h.logger.Error("NewReadLineFromEnd json Unmarshal err", zap.Error(err))
						//ctx.AbortWithError()
					}

					if logsParse.HTTPCode == 0 {
						continue
					}

					if index >= logSize {
						break
					}

					data := logsData{
						Content:     string(line),
						Level:       logsParse.Level,
						Time:        logsParse.Time,
						Path:        logsParse.Path,
						Method:      logsParse.Method,
						Msg:         logsParse.Msg,
						HTTPCode:    logsParse.HTTPCode,
						TraceID:     logsParse.TraceID,
						CostSeconds: logsParse.CostSeconds,
					}

					res.Logs[index] = data
					index++
				}
			} else {
				break
			}
		}
		ctx.Payload(res)
	}
}

// Package httpclient
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 22:08
package httpclient

import (
	"bufio"
	"bytes"

	"go.uber.org/zap"
)

// AlarmVerify 验证解析正文并验证其是否正确
type AlarmVerify func(body []byte) (shouldAlarm bool)

type AlarmObject interface {
	Send(subject, body string) error
}

func onFailedAlarm(title string, raw []byte, logger *zap.Logger, alarmObject AlarmObject) {
	buf := bytes.NewBuffer(nil)

	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		buf.WriteString(scanner.Text())
		buf.WriteString(" <br/>")
	}

	if err := alarmObject.Send(title, buf.String()); err != nil && logger != nil {
		logger.Error("calls failed alarm err", zap.Error(err))
	}
}

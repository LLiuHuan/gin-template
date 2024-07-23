// Package alert
//
//	@program:		gin-template
//	@author:		[lliuhuan](https://github.com/lliuhuan)
//	@create:		2024-07-02 22:40
//	@description:	告警通知
package alert

import (
	"encoding/json"
	"fmt"
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/proposal"
	"github.com/LLiuHuan/gin-template/pkg/errors"
	"github.com/LLiuHuan/gin-template/pkg/httpclient"
	"github.com/LLiuHuan/gin-template/pkg/mail"

	"go.uber.org/zap"
)

// NotifyHandler 告警通知
func NotifyHandler(logger *zap.Logger) func(msg *proposal.AlertMessage) {
	if logger == nil {
		panic("logger required")
	}

	cfg := configs.Get().Notify

	return func(msg *proposal.AlertMessage) {
		switch cfg.Way {
		case configs.ProjectNotifyMail:
			sendEmail(logger, msg)
			break
		case configs.ProjectNotifyWeChat:
			sendWeChat(logger, msg)
			break
		default:
			sendEmail(logger, msg)
		}
		return
	}
}

// sendEmail 发送邮件
func sendEmail(logger *zap.Logger, msg *proposal.AlertMessage) {
	cfg := configs.Get().Notify.Mail
	if cfg.Host == "" || cfg.Port == 0 || cfg.User == "" || cfg.Pass == "" || cfg.To == "" {
		logger.Error("Mail config error")
		return
	}

	subject, body, err := newHTMLEmail(
		msg.Method,
		msg.HOST,
		msg.URI,
		msg.TraceID,
		msg.ErrorMessage,
		msg.ErrorStack,
	)
	if err != nil {
		logger.Error("email template error", zap.Error(err))
		return
	}

	options := &mail.Options{
		MailHost: cfg.Host,
		MailPort: cfg.Port,
		MailUser: cfg.User,
		MailPass: cfg.Pass,
		MailTo:   cfg.To,
		Subject:  subject,
		Body:     body,
	}
	if err := mail.Send(options); err != nil {
		logger.Error("发送告警通知邮件失败", zap.Error(errors.WithStack(err)))
	}
}

// sendWeChat 发送微信
func sendWeChat(logger *zap.Logger, msg *proposal.AlertMessage) {
	cfg := configs.Get().Notify.WeChat
	if cfg.Key == "" {
		logger.Error("WeChat config error")
		return
	}

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", cfg.Key)
	data := json.RawMessage(`{"msgtype": "markdown", "markdown": {"content": "测试消息哦"}}`)
	body, err := httpclient.PostJSON(url, data)
	fmt.Println(body, err)

	return
}

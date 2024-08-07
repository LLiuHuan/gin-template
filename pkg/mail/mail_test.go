// Package mail
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 22:14
package mail

import (
	"testing"
)

func TestSend(t *testing.T) {
	options := &Options{
		MailHost: "smtp.163.com",
		MailPort: 465,
		MailUser: "xxx@163.com",
		MailPass: "", //密码或授权码
		MailTo:   "",
		Subject:  "subject",
		Body:     "body",
	}
	err := Send(options)
	if err != nil {
		t.Error("Mail Send error", err)
		return
	}
	t.Log("success")
}

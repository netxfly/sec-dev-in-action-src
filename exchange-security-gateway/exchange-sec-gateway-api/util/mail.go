package util

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"

	"exchange_zero_trust_api/logger"
	"exchange_zero_trust_api/vars"
)

func SendMail(username, content string) (result bool, err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%v@%v", vars.MailUser, vars.MailSuffix))
	m.SetHeader("To", fmt.Sprintf("%v@%v", username, vars.MailSuffix))
	m.SetHeader("Subject", fmt.Sprintf("[邮箱激活] 您有新的设备需要激活"))
	m.SetBody("text/html", fmt.Sprintf("\n\n%v\n\n", content))

	d := gomail.NewDialer(vars.MailServer, vars.MailPort, vars.MailUser, vars.MailPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		logger.Log.Error(err)
	}
	return result, err
}

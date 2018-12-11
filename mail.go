package main

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

type MailOptions struct {
	From     string
	Host     string
	Port     string
	UserName string
	Password string
}

func (p PostData) toHTML() string {
	return fmt.Sprintf("<a href=\"%v\">%v</a> %v руб.", p.Href, p.Title, p.Price)
}

func genHTML(posts []PostData) []byte {
	var res string
	res = res + "<h2>Новые поступления:<h2>"
	for _, post := range posts {
		res = res + "<br>" + post.toHTML()
	}
	return []byte(res)
}

func sendMails(posts []PostData) {
	opts := getMailOptions()
	mail := email.NewEmail()
	mail.To = []string{opts.From}
	mail.From = opts.From
	mail.HTML = genHTML(posts)
	mail.Subject = "новые позиции!"
	err := mail.SendWithTLS(opts.Host+opts.Port, smtp.PlainAuth("", opts.UserName, opts.Password, opts.Host),
		&tls.Config{
			ServerName: "smtp.rambler.ru",
		})
	if err != nil {
		log.Println("send email error: ", err)
	}
}

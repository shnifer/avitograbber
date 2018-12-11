package main

import (
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

type MailOptions struct {
	From     string
	Addr     string
	UserName string
	Password string
	Host     string
}

func (p PostData) toHTML() string {
	return fmt.Sprintf("<a href=\"%v\">%v</a> $%v руб.", p.Href, p.Title, p.Price)
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
	emails := getEmails()
	if len(emails) == 0 {
		return
	}
	mail := email.NewEmail()
	mail.To = emails
	mail.From = opts.From
	mail.HTML = genHTML(posts)
	err := mail.Send(opts.Addr, smtp.PlainAuth("", opts.UserName, opts.Password, opts.Host))
	if err != nil {
		log.Println("send email error: ", err)
	}
}

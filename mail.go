package main

import (
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

func sendMails(posts []PostData) {
	opts := getMailOptions()
	emails := getEmails()
	if len(emails) == 0 {
		return
	}
	mail := email.NewEmail()
	mail.To = emails
	mail.From = opts.From
	mail.HTML = []byte("Test <B>HTML</B>")
	err := mail.Send(opts.Addr, smtp.PlainAuth("", opts.UserName, opts.Password, opts.Host))
	if err != nil {
		log.Println("send email error: ", err)
	}
}

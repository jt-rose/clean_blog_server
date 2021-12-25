package utils

import (
	"fmt"
	"net/smtp"

	"github.com/jt-rose/clean_blog_server/constants"
)


func SendSampleEmail() {
	// sender data
	from := constants.ENV_VARIABLES.EMAIL_ADDRESS
	password := constants.ENV_VARIABLES.EMAIL_PASSWORD
	// receiver address
	toEmail := constants.ENV_VARIABLES.SAMPLE_TO_EMAIL
	to := []string{toEmail}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// message
	subject := "Subject: Our Golang Email\n"
	body := "our first email!"
	message := []byte(subject + body)
	// athentication data
	// func PlainAuth(identity, username, password, host string) Auth
	auth := smtp.PlainAuth("", from, password, host)
	// send mail
	// func SendMail(addr string, a Auth, from string, to []string, msg []byte) error
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("mail sent")
}
package utils

import (
	"fmt"
	"net/smtp"

	"github.com/jt-rose/clean_blog_server/constants"
)

func SendPasswordResetEmail(recieverEmail string, resetKey string) error {
	// sender data
	from := constants.ENV_VARIABLES.EMAIL_ADDRESS
	password := constants.ENV_VARIABLES.EMAIL_PASSWORD
	// receiver address
	toEmail := recieverEmail
	to := []string{toEmail}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// message
	subject := "Clean Blog Password Reset Link\n"
	body := "A request to reset your password on Clean Blog was recently made. Please visit the following link to reset your password:\n" +
	fmt.Sprintf("%s/reset-password/%s", constants.ENV_VARIABLES.FRONTEND_URL, resetKey)

	message := []byte(subject + body)
	// athentication data
	// func PlainAuth(identity, username, password, host string) Auth
	auth := smtp.PlainAuth("", from, password, host)
	// send mail
	// func SendMail(addr string, a Auth, from string, to []string, msg []byte) error
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		return err
	}
	fmt.Println("password reset requested")
	return nil
}
package mailer

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type MailerService struct{}

func (s *MailerService) SendWarningEmail(email string) {
	message := gomail.NewMessage()

	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "IP address changed")

	message.SetBody("text/plain", "Your IP address has changed. Please, review your account's activity.")

	// I masked the password to not make my credentials public domain
	dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "4636a6e3bed3f0", "**********b47")

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error sending warning email:", err)
	} else {
		fmt.Println("Warning email sent")
	}
}

package service

import (
	"fmt"
	"log"
	"notifier/config"

	"gopkg.in/gomail.v2"
)

func SendEmail(cfg *config.AppConfig, emails []string, body string, subject string) error {

	dialer := gomail.NewDialer(cfg.SmtpHost, 587, cfg.SmtpUser, cfg.SmtpPass)
	server, err := dialer.Dial()
	if err != nil {
		return err
	}
	defer server.Close()

	successCount := 0
	message := gomail.NewMessage()
	for _, email := range emails {
		message.SetHeader("From", cfg.SmtpUser)
		message.SetAddressHeader("To", email, "")
		message.SetHeader("Subject", subject)
		message.SetBody("text/html", body)

		if err = gomail.Send(server, message); err != nil {
			log.Printf("Failed to send email to %s: %v", email, err)
		} else {
			successCount++
			log.Printf("Email sent successfully to %s", email)
		}
		message.Reset() // 다음 이메일을 위해 메시지 초기화
	}

	if successCount == 0 {
		return fmt.Errorf("failed to send email to all subscribers")
	}

	return nil
}

package service

import (
	"log"
	"notifier/config"

	"gopkg.in/gomail.v2"
)

func SendEmail(cfg *config.AppConfig, emailData []map[string]interface{}) error {

	dialer := gomail.NewDialer(cfg.SmtpHost, 587, cfg.SmtpUser, cfg.SmtpPass)
	server, err := dialer.Dial()
	if err != nil {
		return err
	}
	defer server.Close()

	for _, data := range emailData {
		email := data["email"].(string)
		title := data["title"].(string)
		body := data["body"].(string)

		message := gomail.NewMessage()
		message.SetHeader("From", cfg.SmtpUser)
		message.SetAddressHeader("To", email, "")
		message.SetHeader("Subject", title)
		message.SetBody("text/html", body)

		if err = gomail.Send(server, message); err != nil {
			log.Printf("Failed to send email to %s: %v", email, err)
			// TODO: 추후 Fallback 로직 추가 (DLQ를 활용하던가 아님 Fallback 로직을 여기에다가 추가하던가..)
		} else {
			log.Printf("Email sent successfully to %s", email)
		}
	}

	return nil
}

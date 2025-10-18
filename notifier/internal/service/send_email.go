package service

import (
	"fmt"
	"log"
	"notifier/config"

	"gopkg.in/gomail.v2"
)

func SendEmail(cfg *config.AppConfig, emails []string, body string, category string) error {

	dialer := gomail.NewDialer(cfg.SmtpHost, 587, cfg.SmtpUser, cfg.SmtpPass)
	server, err := dialer.Dial()
	if err != nil {
		return err
	}
	defer server.Close()

	successCount := 0
	subject := getEmailSubject(category)
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

func getEmailSubject(category string) string {
	if category == "ssu_path" {
		return "[숭실대학교 IT지원위원회] 신규 SSU_PATH 비교과 프로그램 등록 알림"
	} else {
		return "[숭실대학교 IT지원위원회] 신규 숭실대학교 공지사항 등록 알림"
	}
}

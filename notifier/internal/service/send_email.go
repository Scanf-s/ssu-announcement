package service

import (
	"fmt"
	"log"
	"net/smtp"
	"notifier/internal/dto"

	"notifier/config"
)

func SendEmail(cfg *config.AppConfig, emails []string, announcement dto.Announcement) error {
	// 이메일 발송 로직
	emailBodyTemplate :=
		`
		안녕하세요 숭실대학교 공지사항 알림 서비스입니다.
		구독하신 %s 카테고리의 공지사항이 새로 등록되었습니다.\n\n
		등록일자: %s
		등록부서: %s
		상세링크: %s
		진행상태: %s

		숭실대학교 IT지원위원회 공지사항 알림 서비스를 이용해주셔서 감사드립니다.
		`

	for _, email := range emails {
		subject := "[숭실대학교 공지사항 등록 알림] " + announcement.Title
		body := fmt.Sprintf(
			emailBodyTemplate,
			announcement.Category,
			announcement.Date,
			announcement.Department,
			announcement.Link,
			announcement.Status,
		)

		// 이메일 메시지 헤더
		msg := []byte(
			"From: " + cfg.SmtpUser + "\r\n" +
				"To: " + email + "\r\n" +
				"Subject: " + subject + "\r\n" +
				"\r\n" +
				body + "\r\n",
		)

		err := smtp.SendMail(
			cfg.SmtpHost+":"+cfg.SmtpPort,
			cfg.Auth,
			cfg.SmtpUser,
			[]string{email},
			msg,
		)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", email, err)
			continue
		}
		log.Printf("Email sent successfully to %s", email)
	}

	return nil
}

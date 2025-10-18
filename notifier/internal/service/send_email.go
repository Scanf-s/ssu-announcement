package service

import (
	"context"
	"fmt"
	"log"
	"notifier/internal/repository"

	"notifier/config"
)

func SendEmail(ctx context.Context, cfg *config.AppConfig, data string, category string) error {

	emails, err := repository.GetSubscribers(ctx, cfg, category)
	if err != nil {
		return err
	}
	if len(emails) == 0 {
		return fmt.Errorf("no subscribers for category: %s", category)
	}

	// TODO : 이메일 템플릿을 internal/template에서 embed된 파일 읽어서 사용하도록 구성해야함
	log.Printf("Body : %s", data)
	emailBody, err := CreateEmailTemplate(category, data)
	if err != nil {
		return err
	}

	log.Printf("Email Body : %s", emailBody)
	return nil
}

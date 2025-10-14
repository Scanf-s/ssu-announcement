package service

import (
	"context"
	"log"
	"notifier/internal/repository"

	"notifier/config"
)

func SendEmail(ctx context.Context, cfg *config.AppConfig, body string, category string) error {

	emails, err := repository.GetSubscribers(ctx, cfg, category)
	if err != nil {
		return err
	}
	if len(emails) == 0 {
		log.Printf("No subscribers for %s", category)
		return nil
	}

	// TODO : 이메일 템플릿을 internal/template에서 embed된 파일 읽어서 사용하도록 구성해야함

	return nil
}

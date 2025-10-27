package service

import (
	"context"
	"fmt"
	"notifier/internal/dto"
	"notifier/internal/repository"

	"notifier/config"
)

func NotificationService(ctx context.Context, cfg *config.AppConfig, data dto.Message, category string) error {

	emails, err := repository.GetSubscribers(ctx, cfg, category)
	if err != nil {
		return err
	}
	if len(emails) == 0 {
		return fmt.Errorf("no subscribers for category: %s", category)
	}

	emailBody, err := CreateEmailTemplate(category, data)
	if err != nil {
		return err
	}

	err = SendEmail(cfg, emails, emailBody, data.GetTitle())
	if err != nil {
		return err
	}

	return nil
}

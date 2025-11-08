package service

import (
	"context"
	"fmt"
	"notifier/internal/dto"
	"notifier/internal/repository"

	"notifier/config"
)

func NotificationService(ctx context.Context, cfg *config.AppConfig, data dto.Message, category string) error {

	subscribers, err := repository.GetSubscribers(ctx, cfg, category)
	if err != nil {
		return err
	}
	if len(subscribers) == 0 {
		return fmt.Errorf("no subscribers for category: %s", category)
	}

	emailTemplates, err := CreateEmailTemplate(subscribers, category, data)
	if err != nil {
		return err
	}

	err = SendEmail(cfg, emailTemplates)
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"bytes"
	_ "embed"
	"fmt"
	"notifier/internal/dto"

	"html/template"
)

//go:embed template/announcement_template.html
var announcementTemplate string

//go:embed template/ssu_path_template.html
var ssuPathTemplate string

func CreateEmailTemplate(subscribers []map[string]interface{}, category string, data dto.Message) ([]map[string]interface{}, error) {
	var rawTemplate string
	var baseData interface{}
	var emailBodies []map[string]interface{}

	// Create basic data, determine HTML email raw template
	if category == "ssu_path" {
		rawTemplate = ssuPathTemplate
		ssuPathData, ok := data.(dto.SSUPathMessage)
		if !ok {
			return nil, fmt.Errorf("failed to cast data to SSUPathMessage")
		}
		baseData = ssuPathData
	} else {
		rawTemplate = announcementTemplate
		announcementData, ok := data.(dto.AnnouncementMessage)
		if !ok {
			return nil, fmt.Errorf("failed to cast data to AnnouncementMessage")
		}
		baseData = announcementData
	}

	// Parse(Read) email template
	emailTemplate, err := template.New("EmailTemplate").Parse(rawTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Type assertion and Create email HTML body with data
	if category == "ssu_path" {
		ssuBaseData := baseData.(dto.SSUPathMessage)
		for _, subscriber := range subscribers {
			email := subscriber["Email"].(string)
			token := subscriber["UnsubscribeToken"].(string)

			data := ssuBaseData
			data.Email = email
			data.UnsubscribeToken = token

			emailData, err := createEmailTemplate(emailTemplate, data)
			if err != nil {
				return nil, err
			}
			emailBodies = append(emailBodies, emailData)
		}
	} else {
		announcementBaseData := baseData.(dto.AnnouncementMessage)
		for _, subscriber := range subscribers {
			email := subscriber["Email"].(string)
			token := subscriber["UnsubscribeToken"].(string)

			data := announcementBaseData
			data.Email = email
			data.UnsubscribeToken = token

			emailData, err := createEmailTemplate(emailTemplate, data)
			if err != nil {
				return nil, err
			}
			emailBodies = append(emailBodies, emailData)
		}
	}
	return emailBodies, nil
}

func createEmailTemplate(emailTemplate *template.Template, templateData interface{}) (map[string]interface{}, error) {
	var buffer bytes.Buffer
	err := emailTemplate.Execute(&buffer, templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// templateData를 dto.Message 인터페이스로 타입 단언
	message, ok := templateData.(dto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to cast templateData to dto.Message")
	}

	return map[string]interface{}{
		"email": message.GetEmail(),
		"title": message.GetTitle(),
		"body":  buffer.String(),
	}, nil
}

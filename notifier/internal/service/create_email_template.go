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

func CreateEmailTemplate(category string, data dto.Message) (string, error) {
	var rawTemplate string
	var templateData interface{}

	if category == "ssu_path" {
		rawTemplate = ssuPathTemplate
		ssuPathData, ok := data.(dto.SSUPathMessage)
		if !ok {
			return "", fmt.Errorf("failed to cast data to SSUPathMessage")
		}
		templateData = ssuPathData
	} else { // 다른 카테고리는 숭실대학교 공지사항 데이터로 취급
		rawTemplate = announcementTemplate
		announcementData, ok := data.(dto.AnnouncementMessage)
		if !ok {
			return "", fmt.Errorf("failed to cast data to AnnouncementMessage")
		}
		templateData = announcementData
	}

	tmpl, err := template.New("EmailTemplate").Parse(rawTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buffer bytes.Buffer                   // HTML 내용 담아둘 버퍼
	err = tmpl.Execute(&buffer, templateData) // 템플릿에다가 데이터 적용
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buffer.String(), nil
}

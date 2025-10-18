package service

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"notifier/internal/dto"

	"html/template"
)

//go:embed template/announcement_template.html
var announcementTemplate string

//go:embed template/ssu_path_template.html
var ssuPathTemplate string

func CreateEmailTemplate(category string, data string) (string, error) {
	var rawTemplate string
	var templateData interface{}

	if category == "ssu_path" {
		rawTemplate = ssuPathTemplate
		var ssuPathData dto.SSUPathMessage
		if err := json.Unmarshal([]byte(data), &ssuPathData); err != nil {
			return "", fmt.Errorf("failed to unmarshal SSUPath data: %w", err)
		}
		templateData = ssuPathData
	} else { // 다른 카테고리는 숭실대학교 공지사항 데이터로 취급
		rawTemplate = announcementTemplate
		var announcementData dto.AnnouncementMessage
		if err := json.Unmarshal([]byte(data), &announcementData); err != nil {
			return "", fmt.Errorf("failed to unmarshal Announcement data: %w", err)
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

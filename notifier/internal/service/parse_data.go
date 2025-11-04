package service

import (
	"encoding/json"
	"fmt"
	"notifier/internal/dto"
)

func ParseData(data string, messageType string) (dto.Message, error) {
	if messageType == string(dto.MessageTypeAnnouncement) {
		var announcementData dto.AnnouncementMessage
		if err := json.Unmarshal([]byte(data), &announcementData); err != nil {
			return nil, err
		}
		return announcementData, nil
	} else if messageType == string(dto.MessageTypeSSUPath) {
		var ssuPathData dto.SSUPathMessage
		if err := json.Unmarshal([]byte(data), &ssuPathData); err != nil {
			return nil, err
		}
		return ssuPathData, nil
	} else {
		return nil, fmt.Errorf("invalid message Type")
	}
}

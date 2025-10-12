package dto

type MessageType string

const (
	MessageTypeAnnouncement MessageType = "announcement"
	MessageTypeSSUPath      MessageType = "ssu_path"
)

type Message interface {
	GetLink() string
	GetMessageType() MessageType
}

type AnnouncementMessage struct {
	Link       string
	Category   string
	Title      string
	Date       string
	Department string
	Status     string
}

func (msg AnnouncementMessage) GetLink() string {
	return msg.Link
}

func (msg AnnouncementMessage) GetMessageType() MessageType {
	return MessageTypeAnnouncement
}

type SSUPathMessage struct {
	Title                   string
	Label                   string
	Department              string
	GradePointType          string
	Description             string
	Image                   string
	Link                    string
	ApplicationPeriod       string
	EducationPeriod         string
	ApplicationTarget       string
	ApplicationTargetStatus string
	Mileage                 string
	Applicants              string
	Waitlist                string
	Capacity                string
}

func (msg SSUPathMessage) GetLink() string {
	return msg.Link
}

func (msg SSUPathMessage) GetMessageType() MessageType {
	return MessageTypeSSUPath
}

package dto

type MessageType string

const (
	MessageTypeAnnouncement MessageType = "announcement"
	MessageTypeSSUPath      MessageType = "ssu_path"
)

type Message interface {
	GetLink() string
	GetTitle() string
	GetMessageType() MessageType
}

type AnnouncementMessage struct {
	Link       string // 공지사항 링크
	Category   string // 카테고리
	Title      string // 제목
	Date       string // 작성일
	Department string // 등록부서
	Status     string // 상태
}

func (a AnnouncementMessage) GetLink() string {
	return a.Link
}

func (a AnnouncementMessage) GetTitle() string {
	return a.Title
}

func (a AnnouncementMessage) GetMessageType() MessageType {
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

func (s SSUPathMessage) GetLink() string {
	return s.Link
}

func (s SSUPathMessage) GetTitle() string {
	return s.Title
}

func (s SSUPathMessage) GetMessageType() MessageType {
	return MessageTypeSSUPath
}

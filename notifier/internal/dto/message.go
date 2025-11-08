package dto

type MessageType string

const (
	MessageTypeAnnouncement MessageType = "announcement"
	MessageTypeSSUPath      MessageType = "ssu_path"
)

type Message interface {
	GetLink() string
	GetTitle() string
	GetEmail() string
	GetUnsubscribeToken() string
	GetMessageType() MessageType
}

type AnnouncementMessage struct {
	Link             string // 공지사항 링크
	Category         string // 카테고리
	Title            string // 제목
	Email            string // 구독자 이메일
	UnsubscribeToken string // 구독 해제용 토큰
	Date             string // 작성일
	Department       string // 등록부서
	Status           string // 상태
}

func (a AnnouncementMessage) GetLink() string {
	return a.Link
}

func (a AnnouncementMessage) GetTitle() string {
	return a.Title
}

func (a AnnouncementMessage) GetEmail() string { return a.Email }

func (a AnnouncementMessage) GetUnsubscribeToken() string { return a.UnsubscribeToken }

func (a AnnouncementMessage) GetMessageType() MessageType {
	return MessageTypeAnnouncement
}

type SSUPathMessage struct {
	Title                   string
	Email                   string
	UnsubscribeToken        string
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

func (s SSUPathMessage) GetEmail() string { return s.Email }

func (s SSUPathMessage) GetUnsubscribeToken() string { return s.UnsubscribeToken }

func (s SSUPathMessage) GetMessageType() MessageType {
	return MessageTypeSSUPath
}

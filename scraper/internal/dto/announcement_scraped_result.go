package dto

type AnnouncementScrapedResult struct {
	ScrapedDataType string // 스크래핑된 데이터 타입 (announcement, ssu_path)
	Date            string // 작성일
	Status          string // 상태
	Category        string // 카테고리
	Title           string // 제목
	Department      string // 등록부서
	Link            string // 공지사항 링크
}

func (a AnnouncementScrapedResult) GetTitle() string {
	return a.Title
}

func (a AnnouncementScrapedResult) GetLink() string {
	return a.Link
}

func (a AnnouncementScrapedResult) GetDepartment() string {
	return a.Department
}

func (a AnnouncementScrapedResult) GetScrapedDataType() string {
	return a.ScrapedDataType
}

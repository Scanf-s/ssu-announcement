package dto

type SSUPathScrapedResult struct {
	ScrapedDataType         string // 스크래핑된 데이터 타입 (announcement, ssu_path)
	Title                   string // 프로그램 제목
	Label                   string // 모집 상태 (모집중, 종료)
	Department              string // 등록 부서
	GradePointType          string // 학점 인정 여부
	Description             string // 프로그램 설명
	Image                   string // 프로그램 이미지 URL
	Link                    string // 프로그램 상세 페이지 링크
	ApplicationPeriod       string // 신청 기간
	EducationPeriod         string // 교육 기간
	ApplicationTarget       string // 신청 대상
	ApplicationTargetStatus string // 신청 신분
	Mileage                 string // 마일리지
	Applicants              string // 신청자 수
	Waitlist                string // 대기자 수
	Capacity                string // 모집 정원
}

func (s SSUPathScrapedResult) GetTitle() string {
	return s.Title
}

func (s SSUPathScrapedResult) GetLink() string {
	return s.Link
}

func (s SSUPathScrapedResult) GetDepartment() string {
	return s.Department
}

func (s SSUPathScrapedResult) GetScrapedDataType() string {
	return s.ScrapedDataType
}

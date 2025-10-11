package dto

type ScrapedResult interface {
	GetScrapedDataType() string
	GetTitle() string
	GetLink() string
	GetDepartment() string
}

package main

import (
	"log"

	"ssu-announcement-notifier/config"
	"ssu-announcement-notifier/internal/scraper"
	"ssu-announcement-notifier/internal/service/ssu_announcement_parser"
)

func main() {
	log.Println("ssu-announcement-notifier")

	// 환경변수 불러오기
	cfg := config.LoadConfig()

	// 숭실대학교 공지사항 스크래핑
	var resultHtml []byte
	resultHtml = scraper.ScrapeSSUAnnouncements(cfg)

	// 공지사항 HTML 파싱해서 원하는 정보 추출
	parsedResult := ssu_announcement_parser.ParseSSUAnnouncementsHtml(resultHtml)
	log.Println(parsedResult)
}

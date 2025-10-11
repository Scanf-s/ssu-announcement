package ssu_announcement_parser

import (
	"log"
	"strings"
	"time"

	"scraper/internal/dto"

	"github.com/PuerkitoBio/goquery"
)

// ParseSSUAnnouncementsHtml /**
// 숭실대학교 공지사항 HTML 구조 파싱
// 어짜피 최신 공지만 필요하니까 동적 스크래핑 없이 첫번째 페이지 HTML 구조만 파싱
// @param html []byte HTML 구조
// @return []interface{} 공지사항 리스트
// */
func ParseSSUAnnouncementsHtml(html []byte) ([]dto.AnnouncementScrapedResult, error) {
	// 오늘 날짜 정보
	startDate := time.Now().AddDate(0, 0, -3).Format("2006.01.02")
	endDate := time.Now().Format("2006.01.02")

	// HTML구조에서 공지사항 리스트 -> notice-lists
	// 현재 HTML 구조
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result []dto.AnnouncementScrapedResult
	doc.Find(".notice-lists").Each(func(i int, ulTag *goquery.Selection) {
		ulTag.Find("li").Each(func(i int, liTag *goquery.Selection) {
			className, _ := liTag.Attr("class")
			var announcementData *dto.AnnouncementScrapedResult

			// start 클래스태그, "" 클래스태그의 경우 date 정보를 가져오는 방법이 다름
			if className == "start" {
				date := liTag.Find(".h2.text-info.font-weight-bold").Text()
				if startDate <= date && date <= endDate {
					data := parseDetails(date, liTag)
					announcementData = &data
				}
			} else if className == "" {
				date := liTag.Find(".h2.text-info.font-weight-bold.d-xl-none").Text()
				if startDate <= date && date <= endDate {
					data := parseDetails(date, liTag)
					announcementData = &data
				}
			} else {
				log.Println("Skip unnecessary li tag")
			}

			if announcementData != nil {
				result = append(result, *announcementData)
			}
		})
	})

	return result, nil
}

func parseDetails(date string, liTag *goquery.Selection) dto.AnnouncementScrapedResult {
	status := strings.TrimSpace(liTag.Find(".notice_col2").Find("span").Text())
	category := strings.TrimSpace(liTag.Find(".notice_col3").Find("span").Find(".label.d-inline-blcok.border.pl-3.pr-3.mr-2").Text())
	title := strings.TrimSpace(liTag.Find(".notice_col3").Find("span").Find(".d-inline-blcok.m-pt-5").Text())
	department := strings.TrimSpace(liTag.Find(".notice_col4").Text())
	link := strings.TrimSpace(liTag.Find(".notice_col3").Find("a").AttrOr("href", ""))

	return dto.AnnouncementScrapedResult{
		ScrapedDataType: "announcement",
		Date:            date,
		Status:          status,
		Category:        category,
		Title:           title,
		Department:      department,
		Link:            link,
	}
}

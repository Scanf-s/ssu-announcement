package ssu_path_parser

import (
	"encoding/json"
	"errors"
	"html"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SSUPathHTMLParser(baseUrl string, targetHtml string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(targetHtml))
	if err != nil {
		log.Fatal("Error loading HTML: ", err)
	}
	log.Printf("Document loaded successfully")

	// .lica_wrap 내부의 ul > li 순회 (바로 하위 li 태그들만 선택, 내부에 li태그가 또 있어서...)
	doc.Find(".lica_wrap > ul > li").Each(func(i int, liTag *goquery.Selection) {
		imagePath, _ := liTag.Find(".img_wrap img#repnImg").Attr("src") // 이미지 경로
		imagePath = baseUrl + imagePath

		label := strings.TrimSpace(liTag.Find(".text_wrap .label_box span").Text())            // 상태 라벨 (모집중, 종료)
		major := strings.TrimSpace(liTag.Find(".text_wrap .major_type .first").Text())         // 전공
		gradePointType := strings.TrimSpace(liTag.Find(".text_wrap .major_type .last").Text()) // 학점인정여부
		title := strings.TrimSpace(liTag.Find(".text_wrap a.tit").Text())                      // 프로그램 제목

		// 상세 페이지 링크 추출
		dataParams, _ := liTag.Find(".text_wrap a.tit").Attr("data-params")
		announcementLink, err := extractDetailLink(dataParams, baseUrl)
		if err != nil {
			log.Println("Error extracting detail link: ", err)
		}

		description := strings.TrimSpace(liTag.Find(".text_wrap p.desc.ellipsis").Text()) // 설명 정보

		// InfoWrap 정보 추출
		var applicationPeriod, educationPeriod, applicationTarget, applicationStatus string
		liTag.Find(".info_wrap dl").Each(func(j int, dlTag *goquery.Selection) {
			dtText := strings.TrimSpace(dlTag.Find("dt").Text())
			ddText := strings.TrimSpace(dlTag.Find("dd").Text())

			switch dtText {
			case "신청기간":
				applicationPeriod = ddText
			case "교육기간":
				educationPeriod = ddText
			case "신청대상":
				applicationTarget = ddText
			case "신청신분":
				applicationStatus = ddText
			}
		})

		// 마일리지 및 신청 정보 추출 (근데 이게 필요한지는 모르겠음)
		var mileage, applicants, waitlist, capacity string
		liTag.Find(".etc_cont .rq_desc li dl").Each(func(k int, dlTag *goquery.Selection) {
			dtText := strings.TrimSpace(dlTag.Find("dt").Text())
			ddText := strings.TrimSpace(dlTag.Find("dd").Text())

			switch dtText {
			case "마일리지":
				mileage = ddText
			case "신청자":
				applicants = ddText
			case "대기자":
				waitlist = ddText
			case "모집정원":
				capacity = ddText
			}
		})

		// 결과 출력
		log.Printf("========== Program #%d ==========", i+1)
		log.Printf("Title: %s", title)
		log.Printf("Status: %s", label)
		log.Printf("Major: %s", major)
		log.Printf("Type: %s", gradePointType)
		log.Printf("Description: %s", description)
		log.Printf("Image: %s", imagePath)
		log.Printf("Link: %s", announcementLink)
		log.Printf("Application Period: %s", applicationPeriod)
		log.Printf("Education Period: %s", educationPeriod)
		log.Printf("Target: %s", applicationTarget)
		log.Printf("Status: %s", applicationStatus)
		log.Printf("Mileage: %s", mileage)
		log.Printf("Applicants: %s", applicants)
		log.Printf("Waitlist: %s", waitlist)
		log.Printf("Capacity: %s", capacity)
		log.Println("=====================================")
	})
}

func extractDetailLink(dataParams string, baseUrl string) (string, error) {
	// HTML 엔티티 디코딩 (&#34; -> ")
	unescapedParams := html.UnescapeString(dataParams)

	// JSON 파싱
	var params map[string]string
	if err := json.Unmarshal([]byte(unescapedParams), &params); err == nil {
		if encSeq, ok := params["encSddpbSeq"]; ok {
			// 실제 상세 페이지 URL 구성
			return baseUrl + "/ptfol/imng/icmpNsbjtPgm/findIcmpNsbjtPgmInfo.do?encSddpbSeq=" + encSeq + "&paginationInfo.currentPageNo=1", nil
		}
	}
	return "", errors.New("could not extract detail link")
}

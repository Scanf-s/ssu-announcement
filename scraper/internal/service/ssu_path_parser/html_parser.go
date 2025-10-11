package ssu_path_parser

import (
	"encoding/json"
	"errors"
	"html"
	"log"
	"scraper/internal/dto"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SSUPathHTMLParser(baseUrl string, targetHtml string) ([]dto.SSUPathScrapedResult, error) {
	log.Println("SSUPathHTMLParser called")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(targetHtml))
	if err != nil {
		return nil, err
	}

	// .lica_wrap 내부의 ul > li 순회 (바로 하위 li 태그들만 선택, 내부에 li태그가 또 있어서...)
	var scrapedResults []dto.SSUPathScrapedResult
	doc.Find(".lica_wrap > ul > li").Each(func(i int, liTag *goquery.Selection) {
		imagePath, _ := liTag.Find(".img_wrap img#repnImg").Attr("src")
		imagePath = baseUrl + imagePath

		label := strings.TrimSpace(liTag.Find(".text_wrap .label_box span").Text())
		department := strings.TrimSpace(liTag.Find(".text_wrap .major_type .first").Text())
		gradePointType := strings.TrimSpace(liTag.Find(".text_wrap .major_type .last").Text())
		title := strings.TrimSpace(liTag.Find(".text_wrap a.tit").Text())
		description := strings.TrimSpace(liTag.Find(".text_wrap p.desc.ellipsis").Text()) // 설명 정보

		// 상세 페이지 링크 추출
		dataParams, _ := liTag.Find(".text_wrap a.tit").Attr("data-params")
		announcementLink, err := extractDetailLink(dataParams, baseUrl)
		if err != nil {
			log.Println(err)
		}

		// InfoWrap 정보 추출
		var applicationPeriod, educationPeriod, applicationTarget, applicationTargetStatus string
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
				applicationTargetStatus = ddText
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

		// 데이터 구성 및 추가
		scrapedResults = append(scrapedResults, dto.SSUPathScrapedResult{
			ScrapedDataType:         "ssu_path",
			Title:                   title,
			Label:                   label,
			Department:              department,
			GradePointType:          gradePointType,
			Description:             description,
			Image:                   imagePath,
			Link:                    announcementLink,
			ApplicationPeriod:       applicationPeriod,
			EducationPeriod:         educationPeriod,
			ApplicationTarget:       applicationTarget,
			ApplicationTargetStatus: applicationTargetStatus,
			Mileage:                 mileage,
			Applicants:              applicants,
			Waitlist:                waitlist,
			Capacity:                capacity,
		})
	})

	return scrapedResults, nil
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

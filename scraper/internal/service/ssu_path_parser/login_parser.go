package ssu_path_parser

import (
	"errors"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SSUPathLoginParser(html string) (string, error) {
	log.Println("Parsing SSU-Path Login Page HTML")

	// HTML구조에서 로그인 폼 정보 추출
	htmlReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Fatal(err)
	}

	// 로그인 링크 정보 추출
	var loginLink string
	doc.Find(".tab_cont.box01.is-active").Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "login") {
			loginLink = href
		}
	})

	if loginLink == "" {
		return "", errors.New("login link not found")
	} else {
		return loginLink, nil
	}
}

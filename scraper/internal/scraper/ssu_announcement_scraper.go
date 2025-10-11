package scraper

import (
	"context"
	"io"
	"log"
	"net/http"
	"scraper/internal/repository"
	"scraper/internal/service/ssu_announcement_parser"
	"scraper/internal/service/ssu_path_parser"
	"time"

	"scraper/config"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func ScrapeSSUAnnouncements(ctx context.Context, cfg *config.AppConfig) error {
	log.Println("Scraping SSU Announcements")

	// Request 생성
	request, err := http.NewRequest("GET", cfg.SSUAnnouncementURL, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}

	// Request 보내서 HTML 응답 받기
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Response body 읽기
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// 공지사항 HTML 파싱해서 원하는 정보 추출
	parsedResult, err := ssu_announcement_parser.ParseSSUAnnouncementsHtml(body)
	if err != nil {
		return err
	}

	// DynamoDB에 저장
	err = repository.SaveScrapedData(ctx, cfg, parsedResult)
	if err != nil {
		return err
	}

	return nil
}

func ScrapeSSUPathPrograms(ctx context.Context, cfg *config.AppConfig) error {
	log.Println("Scraping SSU-Path Programs")
	chromeLauncher := cfg.ChromeLauncher

	// 브라우저 실행
	log.Println("Launch Chromium browser")
	url, err := chromeLauncher.Launch()
	if err != nil {
		return err
	}
	log.Println("Get control url")
	browser := rod.New().ControlURL(url)
	if err := browser.Connect(); err != nil {
		return err
	}
	defer func() {
		browser.Close()
		// Chrome launcher cleanup은 자동으로 처리됨
	}()

	// SSU-Path 로그인 페이지 이동
	log.Println("Navigating to SSU-Path main page")
	page, err := browser.Page(proto.TargetCreateTarget{URL: cfg.SSUPathURL})
	if err != nil {
		return err
	}
	if err := page.WaitLoad(); err != nil {
		return err
	}
	time.Sleep(10 * time.Second)

	// 로그인 페이지 HTML 가져오기
	html, err := page.HTML()
	if err != nil {
		return err
	}

	// 로그인 링크 획득
	loginLink, err := ssu_path_parser.SSUPathLoginParser(html)
	if err != nil {
		return err
	}

	// SSO 로그인 페이지로 이동
	log.Println("Navigating to login page")
	if err := page.Navigate(cfg.SSUPathURL + loginLink); err != nil {
		return err
	}
	if err := page.WaitLoad(); err != nil {
		return err
	}
	time.Sleep(10 * time.Second)

	// SSO 로그인 페이지 HTML 가져오기
	html, err = page.HTML()
	if err != nil {
		return err
	}

	// 로그인 폼에 아이디, 비밀번호 입력 (실제학번, 비밀번호 필요) -> 환경변수로 설정해주세요
	useridElem, err := page.Element("#userid")
	if err != nil {
		return err
	}
	if err := useridElem.Input(cfg.SSUPathID); err != nil {
		return err
	}

	pwdElem, err := page.Element("#pwd")
	if err != nil {
		return err
	}
	if err := pwdElem.Input(cfg.SSUPathPW); err != nil {
		return err
	}

	// SSO 로그인 버튼 클릭
	log.Printf("Pressed login button")
	loginBtn, err := page.Element(".btn_login")
	if err != nil {
		return err
	}
	if err := loginBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}
	if err := page.WaitLoad(); err != nil {
		return err
	}
	time.Sleep(10 * time.Second)

	// 로그인 성공 -> 비교과 프로그램 페이지로 직접 이동
	log.Printf("Navigating to non-curricular programs page...")
	programsURL := cfg.SSUPathURL + "/ptfol/imng/icmpNsbjtPgm/findIcmpNsbjtPgmList.do"
	if err := page.Navigate(programsURL); err != nil {
		return err
	}
	if err := page.WaitLoad(); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	// 프로그램 목록이 로드될 때까지 대기
	log.Printf("Waiting for programs list to load...")
	if _, err := page.Element(".lica_wrap"); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	// 비교과 프로그램 목록 HTML 가져오기
	html, err = page.HTML()
	if err != nil {
		return err
	}

	// 비교과 프로그램 데이터 추출
	log.Printf("Parsing SSU-Path Programs Page HTML...")
	scrapedResults, err := ssu_path_parser.SSUPathHTMLParser(cfg.SSUPathURL, html)
	if err != nil {
		return err
	}

	// DynamoDB에 저장
	err = repository.SaveScrapedData(ctx, cfg, scrapedResults)
	if err != nil {
		return err
	}

	return nil
}

package scraper

import (
	"io"
	"log"
	"net/http"
	"os"
	"scraper/internal/service/ssu_path_parser"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"scraper/config"
)

func ScrapeSSUAnnouncements(cfg *config.AppConfig) []byte {
	log.Println("Scraping SSU Announcements")

	// Request 생성
	request, err := http.NewRequest("GET", cfg.SSUAnnouncementURL, nil)
	if err != nil {
		log.Println(err)
	}

	// Request client
	client := &http.Client{}

	// Request 보내기
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	// Response body 읽기
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	return body
}

func ScrapeSSUPathPrograms(cfg *config.AppConfig) {
	log.Println("Scraping SSU-Path Programs")

	// Lambda 환경 감지 및 Chrome 경로 설정
	var l *launcher.Launcher

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda 환경: Chrome 바이너리 경로 지정
		chromePath := "/opt/chrome/chrome" // 컨테이너 내부 파일시스템에 설치된 Chrome 경로
		l = launcher.New().
			Bin(chromePath).
			Headless(true).
			NoSandbox(true). // Lambda에서는 sandbox 모드 비활성화
			Set("disable-gpu").
			Set("disable-dev-shm-usage").
			Set("disable-setuid-sandbox").
			Set("no-first-run").
			Set("no-zygote").
			Set("single-process")
	} else {
		// 로컬 환경은 기본 설정 사용
		l = launcher.New().
			Headless(true)
	}

	// 브라우저 실행
	url := l.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// SSU-Path 로그인 페이지 이동
	page := browser.MustPage(cfg.SSUPathURL)
	page.MustWaitLoad()
	time.Sleep(10 * time.Second)

	// HTML 가져오기
	html, err := page.HTML()
	if err != nil {
		log.Printf("Failed to get HTML: %v", err)
	}

	// 로그인 링크 가져오기
	loginLink, err := ssu_path_parser.SSUPathLoginParser(html)
	if err != nil {
		log.Fatal(err)
	}

	// 로그인 페이지로 이동
	log.Println("Navigating to login page...")
	page = page.MustNavigate(cfg.SSUPathURL + loginLink)
	page.MustWaitLoad()
	time.Sleep(10 * time.Second)

	// 로그인 페이지 HTML 가져오기
	html, err = page.HTML()
	if err != nil {
		log.Printf("Failed to get HTML: %v", err)
	}

	// 로그인 폼에 아이디, 비밀번호 입력 (실제학번, 비밀번호 필요)
	page.MustElement("#userid").MustInput(cfg.SSUPathID)
	page.MustElement("#pwd").MustInput(cfg.SSUPathPW)

	// 로그인 버튼 클릭
	log.Printf("Pressed login button...")
	page.MustElement(".btn_login").MustClick()
	page.MustWaitLoad()
	time.Sleep(10 * time.Second)

	// 비교과 프로그램 페이지로 직접 이동
	log.Printf("Navigating to non-curricular programs page...")
	programsURL := cfg.SSUPathURL + "/ptfol/imng/icmpNsbjtPgm/findIcmpNsbjtPgmList.do"
	page = page.MustNavigate(programsURL)
	page.MustWaitLoad()
	time.Sleep(5 * time.Second)

	// 프로그램 목록이 로드될 때까지 대기
	log.Printf("Waiting for programs list to load...")
	page.MustElement(".lica_wrap")
	time.Sleep(5 * time.Second)

	// 현재 페이지 HTML 가져오기
	html, err = page.HTML()
	if err != nil {
		log.Printf("Failed to get HTML: %v", err)
		return
	}

	log.Printf("Parsing SSU-Path Programs Page HTML...")
	// 비교과 프로그램 데이터 파싱
	ssu_path_parser.SSUPathHTMLParser(cfg.SSUPathURL, html)
}

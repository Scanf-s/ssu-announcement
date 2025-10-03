package scraper

import (
	"io"
	"log"
	"net/http"

	"ssu-announcement-scraper/config"
)

func ScrapeSSUAnnouncements(cfg *config.AppConfig) []byte {
	log.Println("Scraping SSU Announcements")

	// Request 생성
	request, err := http.NewRequest("GET", cfg.SSUAnnouncementURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Request client
	client := &http.Client{}

	// Request 보내기
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Response body 읽기
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

package scrapper

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"scrapper/database"
	"time"

	"github.com/mattn/go-sqlite3"
)

type JijiScrapper struct {
}

func (JijiScrapper) SiteName() string {
	return "Jiji"
}

func (JijiScrapper) Scrape(url string) {
	db := database.NewDB()
	db.Migrate()
	defer db.DB.Close()

	next_url := url
	MAX_PAGES := 100

	//run 100 times to scrape 1-- pages
	i := 1
	for range MAX_PAGES {
		log.Printf("Scrapping page %v of %v", i, MAX_PAGES)
		next_url = scrapeJiji(db, next_url)
		if next_url != "" {
			i++
		}
		return
	}
}

func scrapeJiji(db *database.Database, siteUrl string) string {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	//return a new get request
	req, err := newGetRequest(siteUrl)
	if err != nil {
		log.Fatal(err)
	}
	//set referer based on site being scraped
	req.Header.Set("Referer", "https://jiji.co.ke/search?query=ram+ddr4")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("request error:", err)
		time.Sleep(5 * time.Second)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Println("read error:", err)
		return ""
	}

	var data responseStruct
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("json error:", err)
		return ""
	}

	//commit the data to database
	err = commitToDatabase(db, data)
	if err != nil && !errors.Is(err, sqlite3.ErrConstraint) {
		log.Printf("%v\n", err)
	}
	time.Sleep(10 * time.Second)
	return data.NextUrl
}

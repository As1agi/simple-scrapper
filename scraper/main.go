package main

import scrapper "scrapper/scrappers"

// ONLY FOR SCRAPPING RAM
func main() {
	engine := scrapper.NewScrapperEngine()

	nameUrlMap := map[string]string{
		"Jiji": "https://jiji.co.ke/api_web/v1/listing?query=ram+ddr4&init_page=true&page=1&webp=true",
	}

	for _, siteScrapper := range engine.Scrappers {
		name := siteScrapper.SiteName()
		url := nameUrlMap[name]
		siteScrapper.Scrape(url)
	}
}

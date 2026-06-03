package scrapper

// Scrapper interface for future reference
type Scrapper interface {
	Scrape(string)
	SiteName() string
}

type ScrapEngine struct {
	Scrappers []Scrapper
}

func NewScrapperEngine() *ScrapEngine {
	return &ScrapEngine{
		Scrappers: []Scrapper{
			&JijiScrapper{},
		},
	}
}

type advert struct {
	Title string `json:"title"`
	Price struct {
		Value int `json:"value"`
	} `json:"price_obj"`
	Region     string `json:"region_name"`
	Details    string `json:"slug"`
	Id         int    `json:"id"`
	CategoryId int    `json:"category_id"`
	Url        string `json:"url"`
}

type responseStruct struct {
	AdvertList struct {
		Adverts []advert `json:"adverts"`
	} `json:"adverts_list"`
	NextUrl string `json:"next_url"`
}

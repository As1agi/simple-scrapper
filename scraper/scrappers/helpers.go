package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"scrapper/database"
	"strings"
)

// todo fix
// newGetRequest returns a get request with required headers
func newGetRequest(siteUrl string) (*http.Request, error) {
	req, err := http.NewRequest("GET", siteUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	return req, nil
}

// commit the response struct to database
func commitToDatabase(db *database.Database, data responseStruct) error {
	//loop through the []adverts and commit them to the database
	for _, item := range data.AdvertList.Adverts {
		//check if it has the correct category ID , 284 == Computer Hardware
		//todo make this into a function
		if item.CategoryId == 284 {
			if !strings.Contains(strings.ToLower(item.Title), "ram") {
				continue
			}
			if !strings.Contains(strings.ToLower(item.Title), "touchscreen") {
				continue
			}
			//commit the item to database
			if err := commitOneItemToDatabase(*db, item); err != nil {
				log.Printf(" %v\n", err)
			} //end of commit
		} //end of category ID check
	} //end of for loop
	return nil
}

// commit one advert item to database
func commitOneItemToDatabase(db database.Database, item advert) error {

	//todo later we also get the slug and pass it as a fallback for the details
	brandName, _ := GetBrandName(item.Details)
	ramType, _ := GetRamType(item.Details)
	ramCapacity, _ := GetRamCapacity(item.Details)

	url := "https://jiji.co.ke" + item.Url
	err := db.Insert(item.Id, "Jiji", ramType, ramCapacity, item.Price.Value, item.Region, brandName, url)
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}
	return nil
}

func GetRamType(details string) (string, error) {
	brandNames := []string{
		"ddr4",
		"ddr5",
		"ddr3",
		"ddr2",
	}

	for _, brand := range brandNames {
		exists := strings.Contains(strings.ToLower(details), brand)
		if exists {
			return brand, nil
		}
	}

	return "unknown", fmt.Errorf("%v\n", "Ram Type Not Found")
}

func GetRamCapacity(details string) (string, error) {
	capacities := []string{
		"128gb",
		"64gb",
		"32gb",
		"16gb",
		"8gb",
		"4gb",
		"2gb",
	}

	for _, brand := range capacities {
		exists := strings.Contains(strings.ToLower(details), brand)
		if exists {
			return brand, nil
		}
	}
	return "unknown", fmt.Errorf("%v\n", "Ram capacity Uknown")
}

// GetBrandName extracts and returns the brand name from the details
func GetBrandName(details string) (string, error) {
	brandNames := []string{
		// Major Manufacturers & Global Brands
		"samsung", "hynix", "micron", "corsair", "kingston",
		"gskill", "crucial", "teamgroup", "adata", "patriot",
		"siliconpower", "pny", "mushkin", "geil", "ramaxel",
		"lexar", "transcend", "vcolor", "oloy", "klevv",

		// Niche Gaming, Boutique, & Regional Brands
		"neoforza", "netac", "juhor", "asgard", "gloway",
		"maxsun", "kingspec", "puskill", "walram", "reeinno",
		"faspeed", "kembona", "tcunice", "asint", "goldenfir",
		"yundoo", "xfx", "galax", " colorful", "biostar", "TwinMOS",

		// Industrial, Server, Enterprise & Heritage Brands
		"innodisk", "apacer", "smartmodular", "atp", "supermicro",
		"owc", "timetec", "nemix", "archmemory", "centon",
	}

	for _, brand := range brandNames {
		exists := strings.Contains(strings.ToLower(details), brand)
		if exists {
			return brand, nil
		}
	}
	err := "Brand Not Found"
	return "unknown", fmt.Errorf("%v", err)
}

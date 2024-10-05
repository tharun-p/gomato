package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dnlo/struct2csv"
	"github.com/gocolly/colly"
)

type AddressStruct struct {
	Type            string `json:"@type"`
	StreetAddress   string `json:"streetAddress"`
	AddressLocality string `json:"addressLocality"`
	AddressRegion   string `json:"addressRegion"`
	PostalCode      string `json:"postalCode"`
	AddressCountry  string `json:"addressCountry"`
}

type Geo struct {
	Type      string `json:"@type"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type AggregateRating struct {
	Type        string `json:"@type"`
	RatingValue string `json:"ratingValue"`
	RatingCount int64  `json:"ratingCount"`
	BestRating  int64  `json:"bestRating"`
	WorstRating int64  `json:"worstRating"`
}

type Placejson struct {
	Context         string          `json:"@context"`
	Type            string          `json:"@type"`
	Name            string          `json:"name"`
	Url             string          `json:"url"`
	OpeningHours    string          `json:"openingHours"`
	Hashmap         string          `json:"hashmap"`
	Menu            string          `json:"menu"`
	Address         AddressStruct   `json:"address"`
	Geo             Geo             `json:"geo"`
	Telephone       string          `json:"telephone"`
	PriceRange      string          `json:"priceRange"`
	PaymentAccepted string          `json:"paymentAccepted"`
	Image           string          `json:"image"`
	ServesCuisine   string          `json:"servesCuisine"`
	AggregateRating AggregateRating `json:"aggregateRating"`
}

var visitedMap map[string]string
var forbiddenURL []string

func main() {
	fmt.Println(" Go Scrapper ")

	file, err := os.Create("RestaurentList.csv")
	if err != nil {
		fmt.Println("Error Opening File")
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var RestaurentList []Placejson
	visitedMap = make(map[string]string)

	c := colly.NewCollector(
		colly.AllowedDomains("www.zomato.com", "www.zomato.com", "https://www.zomato.com"),
	)

	forbiddenURL = append(forbiddenURL, "https://www.zomato.com/live", "https://www.zomato.com/who-we-are", "https://blog.zomato.com", "https://www.zomato.com/careers", "https://www.zomato.com/investor-relations", "https://www.zomato.com/report-fraud", "https://blog.zomato.com/press-kit", "https://www.zomato.com/contact", "https://www.zomato.com/zomaland", "https://www.zomato.com/partner-with-us", "https://www.zomato.com/policies/privacy", "https://www.zomato.com/policies/privacy/", "https://www.zomato.com/policies/security/", "https://www.zomato.com/policies/security", "https://www.zomato.com/policies/terms-of-service/", "https://www.zomato.com/policies/terms-of-service")

	forbiddenURLMap := make(map[string]string)

	for _, v := range forbiddenURL {
		forbiddenURLMap[v] = "forbidden"
	}

	c.OnHTML("a", func(e *colly.HTMLElement) {
		_, isForbidden := forbiddenURLMap[e.Attr("href")]
		_, ok := visitedMap[e.Attr("href")]
		if !ok && !isForbidden {
			visitedMap[e.Attr("href")] = "visited"
			err := e.Request.Visit(e.Attr("href"))
			if err != nil {
				return
			}
		}

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("LOOKUP - ", forbiddenURLMap[r.URL.String()])
		_, ok := forbiddenURLMap[r.URL.String()]
		if ok {
			r.Abort()
		}
	})

	c.OnHTML("script", func(e *colly.HTMLElement) {
		var data Placejson
		err := json.Unmarshal([]byte(e.Text), &data)
		if err == nil && data.OpeningHours != "" {
			fmt.Println("Name : ", data.Name, " Place : ", data.Address.AddressRegion)
			RestaurentList = append(RestaurentList, data)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(" Visiting", r.URL)
	})

	// c.MaxDepth = 5 // Update the value depending on the need.
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:130.0) Gecko/20100101 Firefox/130.0"

	err = c.Visit("https://www.zomato.com/chennai/dine-out")
	if err != nil {
		return
	}

	enc := struct2csv.New()
	rows, err := enc.Marshal(RestaurentList)
	if err != nil {

	}

	for _, value := range rows {
		if err := writer.Write(value); err != nil {
			fmt.Println("Error Writing to File")
		}
	}

}

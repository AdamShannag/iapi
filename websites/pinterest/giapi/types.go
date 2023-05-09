package giapi

import (
	"sync"

	"github.com/gocolly/colly"
)

type GoogleImageApi struct {
	collector colly.Collector
	wait      *sync.WaitGroup
	Mutex     *sync.Mutex
	Website   Website
	Data      Data
}

type Website struct {
	WebsiteName string
	GoogleURL   string
	ImageQuery  string
}

type Data struct {
	ImageData []ImageData
	ImageSize ImageSize
}

type ImageSize struct {
	MinWidth  float64
	MinHeight float64
}

func NewGoogleImageApi(websiteName string, size ImageSize) *GoogleImageApi {
	return &GoogleImageApi{
		collector: *colly.NewCollector(),
		Mutex:     &sync.Mutex{},
		wait:      &sync.WaitGroup{},
		Website: Website{
			WebsiteName: websiteName,
			GoogleURL:   "http://www.google.co.in/search?hl=en&q=",
			ImageQuery:  "#main > div > div > div > a",
		},
		Data: Data{
			ImageSize: size,
			ImageData: []ImageData{},
		},
	}
}

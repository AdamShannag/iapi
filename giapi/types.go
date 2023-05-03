package giapi

import (
	"sync"

	"github.com/gocolly/colly"
)

type GoogleImageApi struct {
	collector       colly.Collector
	googleURL       string
	imageXpathQuery string
	wait            *sync.WaitGroup
}

func NewGoogleImageApi() *GoogleImageApi {
	return &GoogleImageApi{
		collector:       *colly.NewCollector(),
		googleURL:       "https://www.google.com.om/search?source=lnms&sa=X&gbv=1&tbm=isch&q=",
		imageXpathQuery: "/html/body/div[3]/table/tbody/tr[*]/td[*]/div/div/div/div/table/tbody/tr[1]/td/a/div/img",
		wait:            &sync.WaitGroup{},
	}
}

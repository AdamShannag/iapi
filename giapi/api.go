package giapi

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	c "github.com/fatih/color"
	"github.com/gocolly/colly"
)

func (i *GoogleImageApi) Search(q string) []string {
	url := fmt.Sprintf("%s%s", i.googleURL, q)
	urls := []string{}

	i.collector.OnRequest(func(r *colly.Request) {
		fmt.Println(c.CyanString("Visting: "), c.BlueString(r.URL.String()))
	})

	i.collector.OnError(func(_ *colly.Response, err error) {
		c.Red("Something went wrong: ", err)
	})

	i.collector.OnResponse(func(r *colly.Response) {
		fmt.Println(c.CyanString("Page visited: "), c.BlueString(r.Request.URL.String()))
	})

	i.collector.OnXML(i.imageXpathQuery, func(e *colly.XMLElement) {
		urls = append(urls, e.Attr("src"))
	})

	i.collector.OnScraped(func(r *colly.Response) {
		fmt.Println(c.BlueString(r.Request.URL.String()), c.GreenString(" scraped!"))
	})

	i.collector.Visit(url)

	return urls
}

func (i *GoogleImageApi) DownloadImage(url string, at, fname string) {
	defer i.wait.Done()

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	s, err := mime.ExtensionsByType(res.Header["Content-Type"][0])

	if err != nil {
		log.Panicln("invalid extenion", err)
	}

	f, err := os.Create(fmt.Sprintf("%s%s%s", at, fname, s[0]))

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	_, err = io.Copy(f, res.Body)

	if err != nil {
		log.Fatal(err)
	}
}

func (i *GoogleImageApi) DownloadImages(urls []string, at, fname string) {
	err := os.MkdirAll(at, 0777)

	if err != nil {
		log.Panic(err)
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " Downloading..."
	s.Start()

	for _, url := range urls {
		i.wait.Add(1)
		go i.DownloadImage(url, at, i.randomBytes(4, fname))
	}

	i.wait.Wait()

	s.Stop()
}

func (i *GoogleImageApi) randomBytes(size int, s string) string {
	buf := make([]byte, size)

	_, err := rand.Read(buf)
	if err != nil {
		log.Fatalf("error while generating random string: %s", err)
	}

	return fmt.Sprintf("%x-%s", buf, s)
}

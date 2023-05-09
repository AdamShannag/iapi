package giapi

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	c "github.com/fatih/color"
	"github.com/gocolly/colly"
)

var scriptUrls []string

func (i *GoogleImageApi) Search(q string) []ImageData {
	var visitedUrls []string

	url := fmt.Sprintf("%s%s", i.Website.GoogleURL, q)

	i.collector.OnRequest(func(r *colly.Request) {
		fmt.Println(c.CyanString("Visting: "), c.BlueString(r.URL.String()))
	})

	i.collector.OnError(func(_ *colly.Response, err error) {
		c.Red("Something went wrong: ", err)
	})

	i.collector.OnResponse(func(r *colly.Response) {
		visitedUrls = append(visitedUrls, r.Request.URL.String())
		fmt.Println(c.CyanString("Page visited: "), c.BlueString(r.Request.URL.String()))
	})

	i.collector.OnHTML(i.Website.ImageQuery, func(element *colly.HTMLElement) {
		href := element.Attr("href")
		url, _ := strings.CutPrefix(href, `/url?q=`)
		if i.urlExclude(url) &&
			strings.Contains(url, i.Website.WebsiteName) &&
			!slices.ContainsFunc(visitedUrls, func(s string) bool { return strings.Contains(url, s) }) {
			i.wait.Add(1)
			go i.get(url)
		}
	})

	i.collector.OnScraped(func(r *colly.Response) {
		fmt.Println(c.BlueString(r.Request.URL.String()), c.GreenString(" scraped!"))
	})

	i.collector.Visit(url)
	i.wait.Wait()
	scriptUrls = []string{}
	return i.Data.ImageData
}

func (i *GoogleImageApi) urlExclude(url string) bool {
	return !(strings.HasPrefix(url, "/") ||
		strings.Contains(url, "wiki") ||
		strings.Contains(url, "google"))
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

func (i *GoogleImageApi) DownloadUrls(fileName string) {
	data, _ := json.MarshalIndent(i.Data.ImageData, "", " ")

	f, _ := os.Create(fmt.Sprintf("%s.json", fileName))
	defer f.Close()

	f.Write(data)

	f.Sync()
}

func (i *GoogleImageApi) get(link string) {

	defer i.wait.Done()

	i.collector.OnHTML(`script#__PWS_DATA__`, func(element *colly.HTMLElement) {
		if !slices.Contains(scriptUrls, element.Request.URL.String()) {
			i.Mutex.Lock()
			scriptUrls = append(scriptUrls, element.Request.URL.String())
			i.Mutex.Unlock()

			jsonParsed, err := gabs.ParseJSON([]byte(element.Text))
			if err != nil {
				log.Println(err)
			}
			contextBytes := jsonParsed.S("props", "context").Bytes()
			var context Context
			json.Unmarshal(contextBytes, &context)
			imageData := ImageData{
				Context: context,
			}
			for _, child := range jsonParsed.S("props", "initialReduxState", "pins").ChildrenMap() {
				var img Image
				json.Unmarshal(child.Bytes(), &img)
				var imgInfo ImageInfo
				imgB := child.S("images", "orig").Bytes()
				json.Unmarshal(imgB, &imgInfo)
				w := imgInfo.Width
				h := imgInfo.Height
				if (w >= i.Data.ImageSize.MinWidth) && (h >= i.Data.ImageSize.MinHeight) {
					img.Info = imgInfo
					imageData.Images = append(imageData.Images, img)
				}
			}
			i.Mutex.Lock()
			i.Data.ImageData = append(i.Data.ImageData, imageData)
			i.Mutex.Unlock()
		}
	})
	i.collector.Visit(link)
}

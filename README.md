# iapi

Scrape and downlaod images from google images

## Sample code
```
func main() {
	// initialize api
	api := iapi.NewGoogleImageApi()

	// search for images
	urls := api.Search("cars")

	// Download images
	api.DownloadImages(urls, "./images/", "g")
}
```

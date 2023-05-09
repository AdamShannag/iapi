# iapi

Scrape and downlaod images from google images

## Sample code
```
func main() {
	// initialize api
	api := iapi.NewGoogleImageApi("pinterest",
		iapi.ImageSize{
			MinWidth:  500.0,
			MinHeight: 1000.0,
		})

	// search for images
	_ = api.Search("messi+pinterest")

	// Download images urls
	api.DownloadUrls("messi")
}
```

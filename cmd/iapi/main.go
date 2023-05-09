package main

import iapi "github.com/AdamShannag/iapi/giapi"

func main() {
	// initialize api
	api := iapi.NewGoogleImageApi()

	// search for images
	urls := api.Search("cars")

	// Download images
	api.DownloadImages(urls, "./images/", "g")
}

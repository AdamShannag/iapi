package giapi

type ImageData struct {
	Context          Context          `json:"context"`
	Images           []Image          `json:"images"`
	RelatedInterests RelatedInterests `json:"seo_related_interests"`
}

type Context struct {
	Country    string `json:"country"`
	CurrentURL string `json:"current_url"`
	Origin     string `json:"origin"`
	Path       string `json:"path"`
}

type Image struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Info  ImageInfo `json:"orig"`
}

type ImageInfo struct {
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
	URL    string  `json:"URL"`
}

type RelatedInterests struct {
	Name string `json:"name"`
}

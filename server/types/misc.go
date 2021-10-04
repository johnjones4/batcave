package types

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type IFrame struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Configuration struct {
	IFrames          []IFrame     `json:"iframes"`
	WeatherLocations []Coordinate `json:"weatherLocations"`
	RSSFeeds         []string     `json:"rssFeeds"`
}

type Response struct {
	IFrames []IFrame  `json:"iframes"`
	Weather []Weather `json:"weather"`
}

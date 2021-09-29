package types

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Config struct {
	LiveFeeds        []string     `json:"liveFeeds"`
	WeatherLocations []Coordinate `json:"weatherLocations"`
}

package core

import "math"

func (c Coordinate) DistanceTo(b Coordinate) float64 {
	r := 6378.137
	dLat := b.Latitude*math.Pi/180 - c.Latitude*math.Pi/180
	dLon := b.Longitude*math.Pi/180 - c.Longitude*math.Pi/180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(c.Latitude*math.Pi/180)*math.Cos(b.Latitude*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	cc := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := r * cc
	return d * 1000
}

func (e FeedbackError) Error() string {
	return e.message
}

func NewFeedbackError(m string) FeedbackError {
	return FeedbackError{m}
}

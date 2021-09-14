package util

import (
	"errors"
	"fmt"
	"hal9000/types"
)

var ErrorDeviceNotFound = errors.New("device not found")
var ErrorDisplayNotFound = errors.New("display not found")
var ErrorJobNotFound = errors.New("job not found")
var ErrorPersonNotFound = errors.New("person not found")
var ErrorSessionNotFound = errors.New("session not found")
var ErrorWeatherForecastNotAvailable = errors.New("weather forecast not available")

func ErrorNoInterfacesAvailable(p types.Person) error {
	return fmt.Errorf("no interfaces ready for %s", p.GetOriginName())
}

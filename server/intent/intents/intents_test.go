package intents

import (
	"context"
	"errors"
	"fmt"
	"main/core"
	"main/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type intentTestCaseIteration struct {
	request      core.Request
	prepareCalls func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response
	metadata     core.IntentMetadata
	error        error
}

type intentTestCase struct {
	intentActor core.IntentActor
	iterations  []intentTestCaseIteration
}

var errorTestError = errors.New("test")

func TestIntents(t *testing.T) {
	ctrl := gomock.NewController(t)
	tuneIn := mocks.NewMockTuneIn(ctrl)
	push := mocks.NewMockPush(ctrl)
	ha := mocks.NewMockHomeAssistant(ctrl)
	now := time.Now().UTC()
	haGroups := []core.HomeAssistantGroup{
		{
			Names:     []string{"front room light"},
			DeviceIds: []string{"a", "b"},
		},
		{
			Names:     []string{"outside light"},
			DeviceIds: []string{"c", "d"},
		},
		{
			Names:     []string{"all lights"},
			DeviceIds: []string{"a", "d"},
		},
	}
	weather := mocks.NewMockWeather(ctrl)
	geocode := mocks.NewMockGeocoder(ctrl)

	cases := []intentTestCase{
		{
			intentActor: &Play{
				TuneIn: tuneIn,
				Push:   push,
			},
			iterations: []intentTestCaseIteration{
				{
					request: core.Request{},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"audioSource": "source name",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						r := core.Response{
							OutboundMessage: core.OutboundMessage{
								Media: core.Media{
									URL:  "test url",
									Type: core.MediaTypeAudioStream,
								},
								Action: core.ActionPlay,
							},
						}
						tuneIn.EXPECT().GetStreamURL(metadata.IntentParseReceiver["audioSource"]).Return(r.OutboundMessage.Media.URL, nil)
						return r
					},
				},
				{
					request: core.Request{},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"audioSource": "source name",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						tuneIn.EXPECT().GetStreamURL(metadata.IntentParseReceiver["audioSource"]).Return("", errorTestError)
						return core.ResponseEmpty
					},
					error: errorTestError,
				},
				{
					request: core.Request{
						Source:   "source str",
						ClientID: "client id",
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"audioSource":           "source name",
							"recurranceDescription": "recurrance string",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						push.EXPECT().SendRecurring(gomock.Any(), request.Source, request.ClientID, metadata.IntentParseReceiver["recurranceDescription"], actor.IntentLabel(), map[string]any{
							"audioSource": metadata.IntentParseReceiver["audioSource"],
						}).Return(nil)
						return core.ResponseEmpty
					},
				},
				{
					request: core.Request{
						Source:   "source str",
						ClientID: "client id",
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"audioSource":           "source name",
							"recurranceDescription": "recurrance string",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						push.EXPECT().SendRecurring(gomock.Any(), request.Source, request.ClientID, metadata.IntentParseReceiver["recurranceDescription"], actor.IntentLabel(), map[string]any{
							"audioSource": metadata.IntentParseReceiver["audioSource"],
						}).Return(errorTestError)
						return core.ResponseEmpty
					},
					error: errorTestError,
				},
			},
		},
		{
			intentActor: &Remind{
				Push: push,
			},
			iterations: []intentTestCaseIteration{
				{
					request: core.Request{
						Source:   "source str",
						ClientID: "client id",
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"date":        now.Format(time.RFC3339Nano),
							"description": "description str",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						push.EXPECT().SendLater(gomock.Any(), now, request.Source, request.ClientID, core.PushMessage{
							OutboundMessage: core.OutboundMessage{
								Message: core.Message{
									Text: fmt.Sprintf("I'm reminding you about \"%s\".", metadata.IntentParseReceiver["description"]),
								},
							},
						}).Return(nil)
						return core.ResponseEmpty
					},
				},
				{
					request: core.Request{
						Source:   "source str",
						ClientID: "client id",
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"date":                  now.Format(time.RFC3339Nano),
							"recurranceDescription": "recurrance",
							"description":           "description str",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						push.EXPECT().SendRecurring(gomock.Any(), request.Source, request.ClientID, metadata.IntentParseReceiver["recurranceDescription"], actor.IntentLabel(), map[string]any{
							"description": metadata.IntentParseReceiver["description"],
						}).Return(nil)
						return core.ResponseEmpty
					},
				},
				{
					request: core.Request{
						Source:   "source str",
						ClientID: "client id",
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"date":                  now.Format(time.RFC3339Nano),
							"recurranceDescription": "recurrance",
							"description":           "description str",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						push.EXPECT().SendRecurring(gomock.Any(), request.Source, request.ClientID, metadata.IntentParseReceiver["recurranceDescription"], actor.IntentLabel(), map[string]any{
							"description": metadata.IntentParseReceiver["description"],
						}).Return(errorTestError)
						return core.ResponseEmpty
					},
					error: errorTestError,
				},
				{
					request: core.Request{
						Source:   "source str",
						ClientID: "client id",
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"date":        now.Format(time.RFC3339Nano),
							"description": "description str",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						push.EXPECT().SendLater(gomock.Any(), now, request.Source, request.ClientID, core.PushMessage{
							OutboundMessage: core.OutboundMessage{
								Message: core.Message{
									Text: fmt.Sprintf("I'm reminding you about \"%s\".", metadata.IntentParseReceiver["description"]),
								},
							},
						}).Return(errorTestError)
						return core.ResponseEmpty
					},
					error: errorTestError,
				},
			},
		},
		{
			intentActor: &ToggleDevice{
				HomeAssistant: ha,
			},
			iterations: []intentTestCaseIteration{
				{
					request: core.Request{
						Message: core.Message{
							Text: "turn the front room lights on",
						},
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"onOff": "on",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						ha.EXPECT().Groups().Return(haGroups)
						for _, id := range haGroups[0].DeviceIds {
							ha.EXPECT().ToggleDeviceState(id, true)
						}
						return core.ResponseEmpty
					},
				},
				{
					request: core.Request{
						Message: core.Message{
							Text: "turn the outside light off",
						},
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"onOff": "off",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						ha.EXPECT().Groups().Return(haGroups)
						for _, id := range haGroups[1].DeviceIds {
							ha.EXPECT().ToggleDeviceState(id, false)
						}
						return core.ResponseEmpty
					},
				},
				{
					request: core.Request{
						Message: core.Message{
							Text: "turn the outside light off",
						},
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"onOff": "off",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						ha.EXPECT().Groups().Return(haGroups)
						ha.EXPECT().ToggleDeviceState(gomock.Any(), false).Return(errorTestError)
						return core.ResponseEmpty
					},
					error: errorTestError,
				},
			},
		},
		{
			intentActor: &Unknown{},
			iterations: []intentTestCaseIteration{
				{
					request: core.Request{},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"answer": "test string",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						return core.Response{
							OutboundMessage: core.OutboundMessage{
								Message: core.Message{
									Text: metadata.IntentParseReceiver["answer"].(string),
								},
							},
						}
					},
				},
			},
		},
		{
			intentActor: &Weather{
				Weather:  weather,
				Geocoder: geocode,
			},
			iterations: []intentTestCaseIteration{
				{
					request: core.Request{},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"date":     "",
							"location": "test location",
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						coord := core.Coordinate{Latitude: 1, Longitude: 2}
						weatherF := core.WeatherForecast{
							Forecast: []core.WeatherForecastPeriod{
								{
									StartTime:        now,
									EndTime:          now,
									DetailedForecast: "weather forecast",
								},
							},
							RadarURL: "some url",
						}
						loc := metadata.IntentParseReceiver["location"]
						geocode.EXPECT().Geocode(loc).Return(coord, nil)
						weather.EXPECT().PredictWeather(coord).Return(weatherF, nil)
						return core.Response{
							OutboundMessage: core.OutboundMessage{
								Message: core.Message{
									Text: fmt.Sprintf("The weather from %s to %s in %s will be: %s",
										weatherF.Forecast[0].StartTime.Format(core.FriendlyDateFormat),
										weatherF.Forecast[0].EndTime.Format(core.FriendlyDateFormat),
										loc,
										weatherF.Forecast[0].DetailedForecast,
									),
								},
								Media: core.Media{
									URL:  weatherF.RadarURL,
									Type: core.MediaTypeImage,
								},
							},
						}
					},
				},
				{
					request: core.Request{
						Coordinate: core.Coordinate{
							Latitude:  1,
							Longitude: 2,
						},
					},
					metadata: core.IntentMetadata{
						IntentParseReceiver: map[string]any{
							"date": now.Add(time.Hour * 6).Format(time.RFC3339Nano),
						},
					},
					prepareCalls: func(request core.Request, actor core.IntentActor, metadata core.IntentMetadata) core.Response {
						weatherF := core.WeatherForecast{
							Forecast: []core.WeatherForecastPeriod{
								{
									StartTime:        now,
									EndTime:          now.Add(time.Hour),
									DetailedForecast: "weather forecast",
								},
								{
									StartTime:        now.Add(time.Hour),
									EndTime:          now.Add(time.Hour * 10),
									DetailedForecast: "weather forecast",
								},
								{
									StartTime:        now.Add(time.Hour * 10),
									EndTime:          now.Add(time.Hour * 20),
									DetailedForecast: "weather forecast",
								},
							},
							RadarURL: "some url",
						}
						weather.EXPECT().PredictWeather(request.Coordinate).Return(weatherF, nil)
						return core.Response{
							OutboundMessage: core.OutboundMessage{
								Message: core.Message{
									Text: fmt.Sprintf("The weather from %s to %s in your location will be: %s",
										weatherF.Forecast[1].StartTime.Format(core.FriendlyDateFormat),
										weatherF.Forecast[1].EndTime.Format(core.FriendlyDateFormat),
										weatherF.Forecast[1].DetailedForecast,
									),
								},
							},
						}
					},
				},
			},
		},
	}

	for _, c := range cases {
		for _, iteration := range c.iterations {
			expectedResponse := iteration.prepareCalls(iteration.request, c.intentActor, iteration.metadata)
			res, err := c.intentActor.ActOnIntent(context.Background(), &iteration.request, &iteration.metadata)
			assert.Equal(t, expectedResponse, res)
			assert.Equal(t, iteration.error, err)
		}
	}
}

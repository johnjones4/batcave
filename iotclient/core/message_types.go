package core

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Message struct {
	Text  string `json:"text"`
	Audio struct {
		Data string `json:"data"`
	} `json:"audio"`
}

type Request struct {
	EventId    string     `json:"eventId"`
	Message    Message    `json:"message"`
	Source     string     `json:"source"`
	ClientID   string     `json:"clientId"`
	Coordinate Coordinate `json:"coordinate"`
}

type Media struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type ResponseBody struct {
	EventId string  `json:"eventId"`
	Message Message `json:"message"`
	Media   Media   `json:"media"`
	Action  string  `json:"action"`
}

type Response struct {
	Type        string        `json:"type"`
	Request     *Request      `json:"request"`
	Response    *ResponseBody `json:"response"`
	PushMessage *ResponseBody `json:"push"`
}

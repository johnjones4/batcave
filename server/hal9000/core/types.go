package core

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type InboundBody struct {
	Body     string     `json:"body"`
	Audio    Audio      `json:"audio"`
	Location Coordinate `json:"location"`
}

type Audio struct {
	Data     string `json:"data"`
	MimeType string `json:"mimeType"`
	GZipped  bool   `json:"gzipped"`
}

type Inbound struct {
	InboundBody
	Command       string `json:"command"`
	State         string `json:"state"`
	User          User   `json:"user"`
	Client        Client `json:"client"`
	ParseMetadata struct {
		Intent string `json:"intent"`
		Body   string `json:"body"`
	} `json:"parseMetadata"`
}

type OutboundBody struct {
	Body  string `json:"body"`
	Media string `json:"media"`
	URL   string `json:"url"`
}

type Outbound struct {
	OutboundBody
	State string `json:"state"`
}

type CommandInfo struct {
	Description  string `json:"description"`
	RequiresBody bool   `json:"requiresBody"`
}

type Intent interface {
	SupportedComandsForState(s string) map[string]CommandInfo
	Execute(req Inbound) (Outbound, error)
}

type FeedbackError struct {
	message string
}

type User struct {
	Name string `json:"name"`
}

type Client struct {
	ID           string   `json:"id"`
	Capabilities []string `json:"capabilities"`
	Users        []string `json:"users"`
}

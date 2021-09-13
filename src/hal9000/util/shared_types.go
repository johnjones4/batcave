package util

type ResponseMessage struct {
	Text  string      `json:"text"`
	URL   string      `json:"url"`
	Extra interface{} `json:"extra"`
}

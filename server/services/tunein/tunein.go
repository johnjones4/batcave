package tunein

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrorStationNotFound = errors.New("the station could not be found")
)

type TuneIn struct {
}

type opml struct {
	Body opmlBody `xml:"body"`
}

type opmlBody struct {
	Outlines []opmlOutline `xml:"outline"`
}

type opmlOutline struct {
	XMLName xml.Name `xml:"outline"`
	URL     string   `xml:"URL,attr"`
	Type    string   `xml:"type,attr"`
}

func (t *TuneIn) GetStreamURL(query string) (string, error) {
	u, err := t.search(query)
	if err != nil {
		return "", err
	}

	res, err := http.Get(u)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (t *TuneIn) search(query string) (string, error) {
	q := make(url.Values)
	q.Set("query", query)
	res, err := http.Get("http://opml.radiotime.com/Search.ashx?" + q.Encode())
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var o opml
	err = xml.Unmarshal(body, &o)
	if err != nil {
		return "", err
	}

	for _, outline := range o.Body.Outlines {
		if outline.Type == "audio" {
			return strings.Trim(outline.URL, " \n"), nil
		}
	}

	return "", ErrorStationNotFound
}

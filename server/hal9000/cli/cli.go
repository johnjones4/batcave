package cli

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/johnjones4/hal-9000/server/hal9000/core"

	"github.com/swaggest/rest"
)

type CLI struct {
	scheme       string
	host         string
	reader       *bufio.Reader
	settingsPath string
	location     core.Coordinate
	settings     struct {
		Key      string `json:"key"`
		ClientId string `json:"clientId"`
	}
}

func New(scheme, host, settingsPath string) *CLI {
	return &CLI{
		scheme:       scheme,
		host:         host,
		settingsPath: settingsPath,
		reader:       bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) Run() {
	err := c.load()
	if err != nil {
		panic(err)
	}
	err = c.ping()
	if err != nil {
		panic(err)
	}
	err = c.discoverLocation()
	if err != nil {
		panic(err)
	}
	for {
		err := c.next()
		if err != nil {
			c.printError(err)
		}
	}
}

func (c *CLI) load() error {
	bytes, err := os.ReadFile(c.settingsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return json.Unmarshal(bytes, &c.settings)
}

func (c *CLI) printError(e error) {
	if er, ok := e.(rest.ErrResponse); ok {
		fmt.Printf("Error number %d: %s\n", er.AppCode, er.Error())
	} else {
		fmt.Printf("Error: %s\n", e.Error())
	}
}

func (c *CLI) ping() error {
	var res map[string]interface{}
	err := c.request("GET", "/api/ping", nil, &res)
	if err != nil {
		return err
	}
	return nil
}

func (c *CLI) next() error {
	input, err := c.prompt("")
	if err != nil {
		return err
	}

	if input == "exit" {
		os.Exit(0)
	}

	req := core.InboundBody{
		Body:     input,
		Location: c.location,
	}
	var res core.OutboundBody
	err = c.request("POST", "/api/request", req, &res)
	if err != nil {
		return err
	}

	fmt.Printf("HAL> %s\n", res.Body)
	if res.URL != "" {
		fmt.Println(res.URL)
	}
	if res.Media != "" {
		fmt.Println(res.Media)
	}

	return nil
}

func (c *CLI) prompt(prompt string) (string, error) {
	fmt.Print(prompt + "> ")
	str, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(str), nil
}

func (c *CLI) request(method string, path string, body interface{}, response interface{}) error {
	var err error
	reqBytes := []byte{}

	if body != nil {
		reqBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	u := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path,
	}

	req, err := http.NewRequest(method, u.String(), io.NopCloser(bytes.NewBuffer(reqBytes)))
	if err != nil {
		return err
	}

	contentType := "application/json"
	reqTime := time.Now().Format(time.RFC3339)

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Request-Time", reqTime)
	req.Header.Set("User-Agent", c.settings.ClientId)

	sigString := strings.Join([]string{c.settings.ClientId, reqTime, contentType}, ":")
	h := hmac.New(sha256.New, []byte(c.settings.Key))
	h.Write([]byte(sigString))
	sha := hex.EncodeToString(h.Sum(nil))
	req.Header.Set("X-Signature", sha)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		var resErr rest.ErrResponse
		err = json.Unmarshal(responseBytes, &resErr)
		if err != nil {
			return err
		}

		return resErr
	} else {
		err = json.Unmarshal(responseBytes, response)
		if err != nil {
			return err
		}

		return nil
	}
}

type ipResponse struct {
	IP string `json:"ip"`
}

type locResponse struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (c *CLI) discoverLocation() error {
	res, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resp ipResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	res, err = http.Get("http://ip-api.com/json/" + resp.IP)
	if err != nil {
		return err
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resp1 locResponse
	err = json.Unmarshal(body, &resp1)
	if err != nil {
		return err
	}

	c.location = core.Coordinate{
		Latitude:  resp1.Latitude,
		Longitude: resp1.Longitude,
	}

	return nil
}

package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/johnjones4/hal-9000/server/hal9000/core"

	"github.com/swaggest/rest"
)

type CLI struct {
	host      string
	tokenPath string
	reader    *bufio.Reader
	token     core.Token
	location  core.Coordinate
}

func New(host string, tokenPath string) *CLI {
	return &CLI{
		host:      host,
		tokenPath: tokenPath,
		reader:    bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) Run() {
	err := c.loadToken()
	if err != nil {
		panic(err)
	}

	err = c.discoverLocation()
	if err != nil {
		panic(err)
	}
	for {
		if c.token.IsExpired() {
			err := c.login()
			if err != nil {
				c.printError(err)
			}
		} else {
			err := c.next()
			if err != nil {
				c.printError(err)
			}
		}
	}
}

func (c *CLI) loadToken() error {
	bytes, err := os.ReadFile(c.tokenPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return json.Unmarshal(bytes, &c.token)
}

func (c *CLI) storeToken() error {
	if c.tokenPath == "" {
		return nil
	}
	bytes, err := json.Marshal(c.token)
	if err != nil {
		return err
	}
	return os.WriteFile(c.tokenPath, bytes, 0600)
}

func (c *CLI) printError(e error) {
	if er, ok := e.(rest.ErrResponse); ok {
		fmt.Printf("Error number %d: %s\n", er.AppCode, er.Error())
	} else {
		fmt.Printf("Error: %s\n", e.Error())
	}
}

func (c *CLI) login() error {
	username, err := c.prompt("Username")
	if err != nil {
		return err
	}

	password, err := c.prompt("Password")
	if err != nil {
		return err
	}

	req := core.LoginRequest{Username: username, Password: password}
	var res core.Token
	err = c.post("/api/login", req, &res)
	if err != nil {
		return err
	}

	c.token = res

	return c.storeToken()
}

func (c *CLI) next() error {
	input, err := c.prompt(c.token.User)
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
	err = c.post("/api/request", req, &res)
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

func (c *CLI) post(path string, body interface{}, response interface{}) error {
	reqBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   path,
	}

	req, err := http.NewRequest("POST", u.String(), io.NopCloser(bytes.NewBuffer(reqBytes)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", c.token.Token)

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

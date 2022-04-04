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

	"github.com/johnjones4/hal-9000/hal9000/core"

	"github.com/swaggest/rest"
)

type CLI struct {
	host   string
	reader *bufio.Reader
	token  core.Token
}

func New(host string) *CLI {
	return &CLI{
		host:   host,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (c *CLI) Run() {
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
	return nil
}

func (c *CLI) next() error {
	input, err := c.prompt(c.token.User)
	if err != nil {
		return err
	}

	if input == "exit" {
		os.Exit(0)
	}

	if len(input) == 0 || input[0] != '/' {
		return errors.New("input not recognized")
	}

	firstSpace := strings.Index(input, " ")
	var command string
	var body string
	if firstSpace < 0 {
		command = input[1:]
		body = ""
	} else {
		command = input[1:firstSpace]
		body = input[firstSpace:]
	}

	req := core.RequestBody{
		Command: command,
		Body:    body,
		Location: core.Coordinate{
			Latitude:  38.804661,
			Longitude: -77.043610,
		},
	}
	var res core.ResponseBody
	err = c.post("/api/request", req, &res)
	if err != nil {
		return err
	}

	fmt.Printf("HAL> %s\n", res.Message)

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

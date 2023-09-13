package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/core"
	"net/http"
)

type Telegram struct {
	Token          string
	ClientRegistry core.ClientRegistry
}

func (t *Telegram) SendToClient(ctx context.Context, clientId string, message core.PushMessage) (bool, error) {
	var chat Chat
	err := t.ClientRegistry.GetClient(ctx, "telegram", clientId, &chat)
	if err != nil {
		return false, err
	}
	if chat.Id == 0 {
		return false, nil
	}
	_, err = t.SendMessage(OutgoingMessage{
		ChatId: chat.Id,
		Message: Message{
			Text: message.Message.Text,
			//TODO media
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *Telegram) RegisterClient(ctx context.Context, msg IncomingMessage) error {
	return t.ClientRegistry.UpsertClient(ctx, "telegram", fmt.Sprint(msg.From.Id), msg.Chat)
}

func (t *Telegram) callMethod(name string, parameters interface{}, response interface{}) error {
	bodyBytes := []byte{}
	var err error

	if parameters != nil {
		bodyBytes, err = json.Marshal(parameters)
		if err != nil {
			return err
		}
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t.Token, name)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %d %s", resp.StatusCode, string(responseBody))
	}

	err = json.Unmarshal(responseBody, response)
	if err != nil {
		return err
	}

	return nil
}

func (t *Telegram) SendMessage(m OutgoingMessage) (SendMessageResponse, error) {
	var resp SendMessageResponse
	err := t.callMethod("sendMessage", m, &resp)
	if err != nil {
		return SendMessageResponse{}, err
	}
	return resp, nil
}

func (t *Telegram) SendReplyKeyboardMarkupMessage(m OutgoingReplyKeyboardMarkupMessage) (SendMessageResponse, error) {
	var resp SendMessageResponse
	err := t.callMethod("sendMessage", m, &resp)
	if err != nil {
		return SendMessageResponse{}, err
	}
	return resp, nil
}

func (t *Telegram) SendReplyKeyboardRemoveMessage(m OutgoingReplyKeyboardRemoveMessage) (SendMessageResponse, error) {
	var resp SendMessageResponse
	err := t.callMethod("sendMessage", m, &resp)
	if err != nil {
		return SendMessageResponse{}, err
	}
	return resp, nil
}

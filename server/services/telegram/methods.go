package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/core"
	"net/http"
)

type TelegramConfiguration struct {
	Token       string `json:"token"`
	SecretToken string `json:"secretToken"`
}

type Telegram struct {
	Configuration  TelegramConfiguration
	ClientRegistry core.ClientRegistry
}

func (t *Telegram) SendOutbound(ctx context.Context, chatId int, message core.OutboundMessage) error {
	_, err := t.SendMessage(OutgoingMessage{
		ChatId: chatId,
		Message: Message{
			Text: message.Message.Text,
		},
	})
	if err != nil {
		return err
	}

	if message.Media.Type == core.MediaTypeImage {
		_, err = t.sendPhoto(OutgoingPhotoMessage{
			ChatId: chatId,
			Photo:  message.Media.URL,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Telegram) SendToClient(ctx context.Context, clientId string, message core.PushMessage) (bool, error) {
	client, err := t.ClientRegistry.Client(ctx, "telegram", clientId, func(client *core.Client, info string) error {
		var chat Chat
		err := json.Unmarshal([]byte(info), &chat)
		if err != nil {
			return err
		}
		client.Info = chat
		return nil
	})
	if err != nil {
		return false, err
	}
	if client.Id == "" {
		return false, nil
	}
	chat, ok := client.Info.(Chat)
	if !ok {
		return false, errors.New("unexpected info")
	}
	if chat.Id == 0 {
		return false, nil
	}
	err = t.SendOutbound(ctx, chat.Id, message.OutboundMessage)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *Telegram) IsClientPermitted(ctx context.Context, r *http.Request, msg IncomingMessage) (bool, error) {
	if r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != t.Configuration.SecretToken {
		return false, errors.New("access denied")
	}

	if msg.Message.Text == "" || msg.Chat.Type != "private" {
		return false, nil
	}

	client, err := t.ClientRegistry.Client(ctx, "telegram", fmt.Sprint(msg.From.Id), nil)
	if err != nil {
		return false, err
	}
	if client.Id == "" {
		return false, fmt.Errorf("user %d not permitted", msg.From.Id)
	}
	return true, nil
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

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t.Configuration.Token, name)

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

func (t *Telegram) sendPhoto(m OutgoingPhotoMessage) (SendMessageResponse, error) {
	var resp SendMessageResponse
	err := t.callMethod("sendPhoto", m, &resp)
	if err != nil {
		return SendMessageResponse{}, err
	}
	return resp, nil
}

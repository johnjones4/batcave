package api

import (
	"fmt"
	"main/core"
	"main/services/telegram"
	"net/http"
)

func (a *apiConcrete) handleTelegramError(w http.ResponseWriter, receiver telegram.Update, err error) {
	a.Log.Error(err)

	_, err = a.Telegram.SendMessage(telegram.OutgoingMessage{
		ChatId: receiver.Message.Chat.Id,
		Message: telegram.Message{
			Text: fmt.Sprintf("Error: \"%s\"", err.Error()),
		},
	})
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (a *apiConcrete) telegramHandler(w http.ResponseWriter, r *http.Request) {
	//TODO user id filter
	//TODO shared secret
	var receiver telegram.Update
	err := a.readJson(r, &receiver)
	if err != nil {
		a.handleTelegramError(w, receiver, err)
		return
	}

	err = a.Telegram.RegisterClient(r.Context(), receiver.Message)
	if err != nil {
		a.handleTelegramError(w, receiver, err)
		return
	}

	if receiver.Message.Text == "" || receiver.Message.Chat.Type != "private" {
		w.WriteHeader(http.StatusOK)
		return
	}

	req := core.Request{
		Message: core.Message{
			Text: receiver.Message.Text,
		},
		Source:   "telegram",
		ClientID: fmt.Sprint(receiver.Message.From.Id),
	}

	res, err := a.coreHandler(r.Context(), req)
	if err != nil {
		a.handleTelegramError(w, receiver, err)
		return
	}

	if res.Message.Text != "" {
		err = a.Telegram.SendOutbound(r.Context(), receiver.Message.Chat.Id, res.OutboundMessage)
		if err != nil {
			a.handleError(w, err, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

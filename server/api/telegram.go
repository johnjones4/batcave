package api

import (
	"fmt"
	"main/core"
	"main/services/telegram"
	"net/http"
)

func (a *API) handleTelegramError(w http.ResponseWriter, receiver telegram.Update, err error) {
	a.Log.Error(err)

	//TODO
	// _, err = a.Telegram.SendMessage(telegram.OutgoingMessage{
	// 	ChatId: receiver.Message.Chat.Id,
	// 	Message: telegram.Message{
	// 		Text: fmt.Sprintf("Error: \"%s\"", err.Error()),
	// 	},
	// })
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (a *API) telegramHandler(w http.ResponseWriter, r *http.Request) {
	var receiver telegram.Update
	err := a.readJson(r, &receiver)
	if err != nil {
		a.handleTelegramError(w, receiver, err)
		return
	}

	ok, err := a.Telegram.IsClientPermitted(r.Context(), r, receiver.Message.From.Id, receiver.Message.Text, receiver.Message.Chat.Type)
	if err != nil {
		a.handleError(w, err, http.StatusUnauthorized)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	req := core.Request{
		EventId: fmt.Sprintf("telegram_%d", receiver.UpdateId),
		Message: core.Message{
			Text: receiver.Message.Text,
		},
		Source:   "telegram",
		ClientID: fmt.Sprint(receiver.Message.From.Id),
	}

	res, err := a.bundledHandler(r.Context(), &req)
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

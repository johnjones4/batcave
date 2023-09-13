package api

import (
	"fmt"
	"main/core"
	"main/services/telegram"
	"net/http"
)

func (a *apiConcrete) telegramHandler(w http.ResponseWriter, r *http.Request) {
	var receiver telegram.Update
	err := a.readJson(r, &receiver)
	if err != nil {
		a.handleError(w, err, http.StatusBadRequest)
		return
	}

	err = a.Telegram.RegisterClient(r.Context(), receiver.Message)
	if err != nil {
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}

	if receiver.Message.Text == "" {
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
		a.handleError(w, err, http.StatusInternalServerError)
		return
	}

	if res.Message.Text != "" {
		_, err = a.Telegram.SendMessage(telegram.OutgoingMessage{
			ChatId: receiver.Message.Chat.Id,
			Message: telegram.Message{
				Text: res.Message.Text,
				//TODO media
			},
		})
		if err != nil {
			a.handleError(w, err, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

package telegrambot

import (
	"errors"

	"github.com/shoshtari/paroo/internal/pkg"
)

func (t TelegramBotImp) DeleteMessage(chatID, messageID int) error {

	req := make(map[string]int)
	req["chat_id"] = chatID
	req["message_id"] = messageID

	res := make(map[string]bool)

	err := pkg.SendHTTPRequest(t.httpClient, t.getUrl("deleteMessage"), req, &res)
	if err != nil {
		return err
	}

	if !res["ok"] {
		return errors.New("response is not ok")
	}
	return nil
}

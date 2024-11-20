package telegrambot

import (
	"errors"
	"time"

	"github.com/shoshtari/paroo/internal/pkg"
)

type TelegramUpdate struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	} `json:"from"`

	Chat struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"chat"`
	Date     int    `json:"date"`
	EditDate int    `json:"edit_date"`
	Text     string `json:"text"`
}
type telegramGetUpdateResponse struct {
	Ok     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

func (t TelegramBotImp) getUpdates() (<-chan TelegramUpdate, <-chan error) {
	updateChan := make(chan TelegramUpdate)
	errChan := make(chan error)
	go func() {
		ticker := time.NewTicker(time.Millisecond)
		for range ticker.C {
			var res telegramGetUpdateResponse
			err := pkg.SendHTTPRequest(t.httpClient, t.getUrl("getUpdates"), nil, &res)
			if err != nil {
				errChan <- err
			}
			if !res.Ok {
				errChan <- errors.New("res is not ok")
			}
			for _, update := range res.Result {
				updateChan <- update
			}

		}
	}()
	return updateChan, errChan
}

func (t TelegramBotImp) runWebhook() (<-chan TelegramUpdate, <-chan error) {
	panic("not implemented")
}

func (t TelegramBotImp) GetUpdatesChan(method string) (<-chan TelegramUpdate, <-chan error) {
	switch method {
	case "getUpdates":
		return t.getUpdates()
	case "webhook":
		return t.runWebhook()
	default:
		var errChan chan error
		go func() {
			errChan <- pkg.UnknownMethodError
		}()
		return nil, errChan
	}

}

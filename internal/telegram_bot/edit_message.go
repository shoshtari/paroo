package telegrambot

import (
	"errors"

	"github.com/shoshtari/paroo/internal/pkg"
)

type EditMessageRequest struct {
	*SendMessageRequest
	MessageID int `json:"message_id"`
}

func NewEditMessageRequest(chatID int, messageID int, text string) EditMessageRequest {
	return EditMessageRequest{
		SendMessageRequest: &SendMessageRequest{
			ChatID: chatID,
			Text:   text,
		},
		MessageID: messageID,
	}

}

func (t TelegramBotImp) EditMessage(request EditMessageRequest) error {
	if request.ChatID == 0 || request.Text == "" {
		return pkg.BadRequestError
	}
	var res SendMessageResponse
	err := pkg.SendHTTPRequest(t.httpClient, t.getUrl("editMessageText"), request, &res)
	if err != nil {
		return err
	}

	if !res.Ok {
		return errors.New("response is not ok")
	}
	return nil
}

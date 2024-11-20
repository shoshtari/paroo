package telegrambot

import (
	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/pkg"
)

type InlineKeybordButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type ReplyMarkup struct {
	ReplyKeybord  [][]string              `json:"reply_keyboard,omitempty"`
	InlineKeybord [][]InlineKeybordButton ` json:"inline_keyboard,omitempty"`
}

type SendMessageRequest struct {
	ChatID      int          `json:"chat_id"`
	Text        string       `json:"text"`
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`
}

type SendMessageResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
}

func NewSendMessageRequest(chatID int, text string) SendMessageRequest {
	return SendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}
}

func WithInlineKeybord(request SendMessageRequest, inlineKeybord [][]InlineKeybordButton) SendMessageRequest {
	request.ReplyMarkup.ReplyKeybord = nil
	request.ReplyMarkup.InlineKeybord = inlineKeybord
	return request
}

func WithReplyKeybord(request SendMessageRequest, replyKeybord [][]string) SendMessageRequest {
	request.ReplyMarkup.InlineKeybord = nil
	request.ReplyMarkup.ReplyKeybord = replyKeybord
	return request
}

func (t TelegramBotImp) SendMessage(request SendMessageRequest) (int, error) {
	if request.ChatID == 0 || request.Text == "" {
		return -1, errors.Wrap(pkg.BadRequestError, "chat id and text is needed")
	}
	var res SendMessageResponse
	err := pkg.SendHTTPRequest(t.httpClient, t.getUrl("sendMessage"), request, &res)
	if err != nil {
		return -1, err
	}

	if !res.Ok {
		return -1, errors.New("response is not ok")
	}
	return res.Result.MessageID, nil

}

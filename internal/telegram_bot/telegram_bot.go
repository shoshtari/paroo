package telegrambot

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/pkg"
)

type TelegramBot interface {
	SendMessage(SendMessageRequest) (int, error)
	EditMessage(EditMessageRequest) error
	DeleteMessage(int, int) error
}

type TelegramBotImp struct {
	httpClient         http.Client
	baseAddress, token string
}

func (t TelegramBotImp) getUrl(path string) string {
	return fmt.Sprintf("%v/bot%v/%v", t.baseAddress, t.token, path)
}

// sendRequest will send the request to telegram bot server.
// if body is nill, it will send a GET else it will be a POST
// body and resbody can be structs of any kind, function will encode/decode json itself
func (t TelegramBotImp) getMe() error {
	return pkg.SendHTTPRequest(t.httpClient, t.getUrl("getMe"), nil, nil)
}

func NewTelegramBot(config configs.SectionTelegram) (TelegramBot, error) {
	var ans TelegramBotImp

	ans.httpClient.Timeout = config.Timeout
	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't parse proxy url")
		}
		ans.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	ans.baseAddress = config.BaseAddress
	ans.token = config.Token

	return ans, ans.getMe()
}

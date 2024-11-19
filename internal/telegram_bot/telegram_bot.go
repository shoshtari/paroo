package telegrambot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/configs"
)

type TelegramBot interface {
	SendMessage(int, string) (int, error)
	EditMessage(int, int, string) error
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
func (t TelegramBotImp) sendRequest(url string, body any, resbody any) error {
	method := http.MethodPost
	if body == nil {
		method = http.MethodGet
	}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "couldn't json marshal request")
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(encodedBody))
	if err != nil {
		return errors.Wrap(err, "couldn't create request")
	}

	res, err := t.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "couldn't send request")
	}

	if res.Body != nil && resbody != nil {
		err = json.NewDecoder(res.Body).Decode(resbody)
		if err != nil {
			return errors.Wrap(err, "couldn't decode response")
		}
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprint("status is ", res.StatusCode))
	}

	return nil

}

func (t TelegramBotImp) getMe() error {
	return t.sendRequest(t.getUrl("getMe"), nil, nil)
}
func (TelegramBotImp) SendMessage(chatID int, text string) (int, error) {
	panic("unimplemented")
}

func (TelegramBotImp) EditMessage(int, int, string) error {
	panic("unimplemented")
}

func (TelegramBotImp) DeleteMessage(int, int) error {
	panic("unimplemented")
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

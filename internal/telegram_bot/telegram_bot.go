package telegrambot

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/pkg"
	"go.uber.org/zap"
)

type TelegramBot interface {
	SendMessage(SendMessageRequest) (int, error)
	EditMessage(EditMessageRequest) error
	DeleteMessage(int, int) error
	GetUpdatesChan(method string) (<-chan TelegramUpdate, <-chan error)
}

type TelegramBotImp struct {
	httpClient           http.Client
	baseAddress, token   string
	lastRecievedUpdateID int
	getUpdateTimeout     int
	logger               *zap.Logger
}

func (t TelegramBotImp) getUrl(path string) string {
	return fmt.Sprintf("%v/bot%v/%v", t.baseAddress, t.token, path)
}

func (t TelegramBotImp) getMe() error {
	return pkg.SendHTTPRequest(t.httpClient, t.getUrl("getMe"), nil, nil)
}

func NewTelegramBot(config configs.SectionTelegram, logger *zap.Logger) (TelegramBot, error) {
	var ans TelegramBotImp

	ans.httpClient.Timeout = config.Timeout
	ans.baseAddress = config.BaseAddress
	ans.token = config.Token
	ans.getUpdateTimeout = config.GetUpdateTimeout
	ans.logger = logger

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't parse proxy url")
		}
		ans.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	return ans, ans.getMe()
}

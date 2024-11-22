package core

import (
	"fmt"

	"github.com/shoshtari/paroo/internal/exchange"
	telegrambot "github.com/shoshtari/paroo/internal/telegram_bot"
)

type ParooCore interface {
	Start() error
}

type ParooCoreImp struct {
	tgbot        telegrambot.TelegramBot
	wallexClient exchange.Exchange
}

func (p ParooCoreImp) Start() error {
	updateChan, errChan := p.tgbot.GetUpdatesChan("getUpdates")
	for {
		select {
		case err := <-errChan:
			return err
		case update := <-updateChan:
			fmt.Println(update)
		}
	}
}

func NewParooCode(tgbot telegrambot.TelegramBot, wallexClient exchange.Exchange) ParooCore {
	return ParooCoreImp{
		tgbot:        tgbot,
		wallexClient: wallexClient,
	}
}

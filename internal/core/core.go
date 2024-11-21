package core

import (
	"fmt"

	telegrambot "github.com/shoshtari/paroo/internal/telegram_bot"
)

type ParooCore interface {
	Start() error
}

type ParooCoreImp struct {
	tgbot telegrambot.TelegramBot
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

func NewParooCode(tgbot telegrambot.TelegramBot) ParooCore {
	return ParooCoreImp{
		tgbot: tgbot,
	}
}

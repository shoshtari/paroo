package core

import (
	"github.com/shoshtari/paroo/internal/pkg"
	telegrambot "github.com/shoshtari/paroo/internal/telegram_bot"
	"go.uber.org/zap"
)

func (p ParooCoreImp) handleBalanceChart(update telegrambot.TelegramUpdate) error {
	balance, err := p.wallexClient.GetTotalBalance()
	if err != nil {
		pkg.GetLogger().Error("couldn't get balance", zap.Error(err))
		return err
	}
	if _, err := p.tgbot.SendMessage(telegrambot.NewSendMessageRequest(update.Message.Chat.ID, balance.String())); err != nil {
		pkg.GetLogger().Error("couldn't send message", zap.Error(err))
	}
	return nil

}

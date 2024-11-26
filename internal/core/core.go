package core

import (
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
	telegrambot "github.com/shoshtari/paroo/internal/telegram_bot"
	"go.uber.org/zap"
)

type ParooCore interface {
	Start() error
}

type UpdateHandler struct {
	Name    string
	Handler func(telegrambot.TelegramUpdate) error
}
type ParooCoreImp struct {
	tgbot        telegrambot.TelegramBot
	wallexClient exchange.Exchange
	balanceRepo  repositories.BalanceRepo
	statRepo     repositories.MarketStatsRepo

	handlers   [][]UpdateHandler
	handlerMap map[string]UpdateHandler
}

func (p ParooCoreImp) handleTelegramNewMessage(update telegrambot.TelegramUpdate) error {
	var req telegrambot.SendMessageRequest

	if handler, exists := p.handlerMap[update.Message.Text]; exists {
		return handler.Handler(update)
	}

	switch update.Message.Text {
	case "/start":
		var replyKeyboard [][]string
		for _, arr := range p.handlers {
			var texts []string
			for _, handler := range arr {
				texts = append(texts, handler.Name)
			}
			replyKeyboard = append(replyKeyboard, texts)
		}

		req = telegrambot.NewSendMessageRequest(update.Message.Chat.ID, "Welcome!")
		req.WithReplyKeybord(replyKeyboard)

	default:
		req = telegrambot.NewSendMessageRequest(update.Message.Chat.ID, "don't know how handle this")
	}

	if _, err := p.tgbot.SendMessage(req); err != nil {
		pkg.GetLogger().Error("couldn't send telegram message", zap.Error(err))
	}

	return nil
}

func (p ParooCoreImp) Start() error {
	updateChan, errChan := p.tgbot.GetUpdatesChan("getUpdates")
	for {
		select {
		case err := <-errChan:
			return err
		case update := <-updateChan:
			if err := p.handleTelegramNewMessage(update); err != nil {
				return err
			}
		}
	}
}

func NewParooCode(tgbot telegrambot.TelegramBot, wallexClient exchange.Exchange, balanceRepo repositories.BalanceRepo, statsRepo repositories.MarketStatsRepo) ParooCore {
	ans := ParooCoreImp{
		tgbot:        tgbot,
		wallexClient: wallexClient,
		balanceRepo:  balanceRepo,
		statRepo:     statsRepo,
	}
	handlers := [][]UpdateHandler{{
		UpdateHandler{"Balance Chart", ans.handleBalanceChart},
	}}

	ans.handlers = handlers
	ans.handlerMap = make(map[string]UpdateHandler)

	for _, row := range handlers {
		for _, handler := range row {
			ans.handlerMap[handler.Name] = handler
		}
	}
	go func() {
		if err := ans.getStatDaemon(); err != nil {
			panic(err)
		}
	}()

	return ans
}

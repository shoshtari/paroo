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
	priceManager PriceManager

	wallexClient exchange.Exchange
	balanceRepo  repositories.BalanceRepo
	marketsRepo  repositories.MarketRepo
	statRepo     repositories.MarketStatsRepo

	handlers   [][]UpdateHandler
	handlerMap map[string]UpdateHandler
}

func (p ParooCoreImp) handleTelegramNewMessage(update telegrambot.TelegramUpdate) error {
	var req telegrambot.SendMessageRequest

	if handler, exists := p.handlerMap[update.Message.Text]; exists {
		return handler.Handler(update)
	}

	pkg.GetLogger().With(
		zap.String("module", "core"),
		zap.String("method", "Handle Telegram message"),
		zap.Int("chat id", update.Message.Chat.ID),
		zap.Int("message id ", update.UpdateID),
	).Debug("new update received")
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

func NewParooCore(
	tgbot telegrambot.TelegramBot, wallexClient exchange.Exchange,
	balanceRepo repositories.BalanceRepo, marketRepo repositories.MarketRepo, statsRepo repositories.MarketStatsRepo,
	priceManager PriceManager,
) ParooCore {
	ans := ParooCoreImp{
		tgbot:        tgbot,
		priceManager: priceManager,
		wallexClient: wallexClient,
		marketsRepo:  marketRepo,
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
	go ans.getStatDaemon()

	return ans
}

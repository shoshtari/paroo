package core

import (
	"context"

	"github.com/pkg/errors"
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

	exchanges []exchange.Exchange

	balanceRepo  repositories.BalanceRepo
	exchangeRepo repositories.ExchangeRepo
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
	tgbot telegrambot.TelegramBot, exchanges []exchange.Exchange,
	balanceRepo repositories.BalanceRepo, marketRepo repositories.MarketRepo, statsRepo repositories.MarketStatsRepo,
	priceManager PriceManager, exchageRepo repositories.ExchangeRepo,
) (ParooCore, error) {
	if len(exchanges) == 0 {
		return nil, errors.Wrap(pkg.BadRequestError, "exchanges slice is empty")
	}

	for _, exchange := range exchanges {
		if err := exchageRepo.Insert(context.TODO(), exchange.GetExchangeInfo()); err != nil {
			pkg.GetLogger("core").Error("couldn't insert exchange to db", zap.Error(err))
			return nil, pkg.InternalError
		}
	}
	ans := ParooCoreImp{
		tgbot:        tgbot,
		priceManager: priceManager,
		exchanges:    exchanges,
		marketsRepo:  marketRepo,
		balanceRepo:  balanceRepo,
		statRepo:     statsRepo,
		exchangeRepo: exchageRepo,
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

	return ans, nil
}

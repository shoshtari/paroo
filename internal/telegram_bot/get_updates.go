package telegrambot

type TelegramUpdate struct{}

func (t TelegramBotImp) GetUpdates(method string) (chan TelegramUpdate, error) {
	panic("not implemented yet!")
}

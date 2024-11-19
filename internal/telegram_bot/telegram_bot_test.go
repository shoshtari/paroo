package telegrambot

import (
	"os"
	"testing"

	"github.com/shoshtari/paroo/test"
	"github.com/stretchr/testify/assert"
)

// since the telegram bot has side effect and IO calls, it's runs are disabled by default
func TestMain(m *testing.M) {
	if os.Getenv("ALL_TESTS") == "" {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestTelegramBot(t *testing.T) {
	config := test.GetTestConfig(t)
	telegramBot, err := NewTelegramBot(config.Telegram)
	assert.Nil(t, err)

	messageID, err := telegramBot.SendMessage(NewSendMessageRequest(config.Telegram.ChatID, "salam"))
	assert.Nil(t, err)
	assert.NotZero(t, messageID)

	// err = telegramBot.EditMessage(config.Telegram.ChatID, messageID, "test2")
	// assert.Nil(t, err)
	//
	// err = telegramBot.DeleteMessage(config.Telegram.ChatID, messageID)
	// assert.Nil(t, err)
}

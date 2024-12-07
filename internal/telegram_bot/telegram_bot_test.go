package telegrambot

import (
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

// since the telegram bot has side effect and IO calls, it's runs are disabled by default
var telegramBot TelegramBot
var config configs.ParooConfig

func TestMain(m *testing.M) {
	if os.Getenv("ALL_TESTS") == "" {
		os.Exit(0)
	}
	var err error
	config = test.GetTestConfig()
	telegramBot, err = NewTelegramBot(config.Telegram, pkg.GetLogger())
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestTelegramBot(t *testing.T) {

	req := NewSendMessageRequest(config.Telegram.ChatID, "salam")
	messageID, err := telegramBot.SendMessage(req)
	assert.Nil(t, err)
	assert.NotZero(t, messageID)

	req2 := EditMessageRequest{
		SendMessageRequest: &req,
		MessageID:          messageID,
	}
	req2.Text = "salam2"

	err = telegramBot.EditMessage(req2)
	assert.Nil(t, err)

	err = telegramBot.DeleteMessage(req.ChatID, messageID)
	assert.Nil(t, err)
}

func TestGetUpdates(t *testing.T) {
	updateChan, errChan := telegramBot.GetUpdatesChan("getUpdates")
	var wg errgroup.Group
	wg.Go(func() error {
		select {
		case err := <-errChan:
			return errors.Wrap(err, "error is not nil")

		case <-time.After(time.Millisecond * 300):
			return nil

		}
	})

	wg.Go(func() error {
		select {
		case <-updateChan:
			return nil
		case <-time.After(time.Millisecond * 300):
			return errors.New("couldn't get update")
		}
	})

	err := wg.Wait()
	assert.Nil(t, err)

}

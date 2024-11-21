/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/core"
	"github.com/shoshtari/paroo/internal/pkg"
	telegrambot "github.com/shoshtari/paroo/internal/telegram_bot"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// runtgbotCmd represents the runtgbot command
var runtgbotCmd = &cobra.Command{
	Use:   "runtgbot",
	Short: "run telegram bot",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := configs.GetConfig()
		if err != nil {
			log.Fatal("couldn't get config, err: ", err)
		}
		logger, err := pkg.GetLogger(config.Log)
		if err != nil {
			log.Fatal("couldn't initialize logger, err: ", err)
		}

		tgbot, err := telegrambot.NewTelegramBot(config.Telegram)
		if err != nil {
			logger.Panic("couldn't initialize telegram bot", zap.Error(err))
		}
		parooCore := core.NewParooCode(tgbot)
		if err := parooCore.Start(); err != nil {
			logger.Fatal("error on running paroo", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(runtgbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runtgbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runtgbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

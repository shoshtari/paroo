/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/shoshtari/paroo/internal/core"
	"github.com/shoshtari/paroo/internal/exchange/wallex"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
	postgresRepo "github.com/shoshtari/paroo/internal/repositories/postgres"

	sqliteRepo "github.com/shoshtari/paroo/internal/repositories/sqlite"
	telegrambot "github.com/shoshtari/paroo/internal/telegram_bot"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func GetRepos(config configs.SectionDatabase) (marketRepo repositories.MarketRepo, err error) {

	if config.Provider == "" {
		if config.Postgres.Host != "" {
			config.Provider = "postgres"
		} else {
			config.Provider = "sqlite"
		}
	}

	switch config.Provider {
	case "postgres":
		ctx := context.Background()
		pgconn, err2 := postgresRepo.ConnectPostgres(ctx, config.Postgres)
		if err != nil {
			err = errors.Wrap(err2, "couldn't connect to postgres")
			return
		}

		marketRepo, err2 = postgresRepo.NewMarketRepo(pgconn, ctx)
		if err2 != nil {
			err = errors.Wrap(err, "couldn't make markets repo")
			return
		}
		return

	case "sqlite":
		marketRepo, err = sqliteRepo.NewMarketRepo(config.Sqlite)
		return
	}
	return marketRepo, err

}

// runtgbotCmd represents the runtgbot command
var runtgbotCmd = &cobra.Command{
	Use:   "runtgbot",
	Short: "run telegram bot",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := configs.GetConfig()
		if err != nil {
			log.Fatal("couldn't get config, err: ", err)
		}

		err = pkg.InitializeLogger(config.Log)
		if err != nil {
			log.Fatal("couldn't initialize logger, err: ", err)
		}

		logger := pkg.GetLogger()
		marketsRepo, err := GetRepos(config.Database)
		if err != nil {
			logger.Fatal("couldn't get repos", zap.Error(err))
		}

		wallexClient, err := wallex.NewWallexClient(config.Wallex, marketsRepo)
		if err != nil {
			logger.Fatal("couldn't connect to wallex", zap.Error(err))
		}

		tgbot, err := telegrambot.NewTelegramBot(config.Telegram)
		if err != nil {
			logger.Panic("couldn't initialize telegram bot", zap.Error(err))
		}
		parooCore := core.NewParooCode(tgbot, wallexClient)

		logger.Info("All dependencies initialized, starting the core")
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

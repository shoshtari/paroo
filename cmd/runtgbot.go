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
	"github.com/shoshtari/paroo/internal/exchange"
	"github.com/shoshtari/paroo/internal/exchange/ramzinex"
	"github.com/shoshtari/paroo/internal/exchange/wallex"
	"github.com/shoshtari/paroo/internal/pkg"
	"github.com/shoshtari/paroo/internal/repositories"
	postgresRepo "github.com/shoshtari/paroo/internal/repositories/postgres"

	telegrambot "github.com/shoshtari/paroo/internal/telegram_bot"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func getRepos(config configs.SectionDatabase) (
	marketRepo repositories.MarketRepo,
	balanceRepo repositories.BalanceRepo,
	statsRepo repositories.MarketStatsRepo,
	exchangeRepo repositories.ExchangeRepo,
	err error) {

	ctx := context.Background()
	pkg.GetLogger().Info("using postgres for db")
	pgconn, err2 := postgresRepo.ConnectPostgres(ctx, config.Postgres)
	if err2 != nil {
		err = errors.Wrap(err2, "couldn't connect to postgres")
		return
	}

	marketRepo, err2 = postgresRepo.NewMarketRepo(ctx, pgconn)
	if err2 != nil {
		err = errors.Wrap(err2, "couldn't make markets repo")
		return
	}

	balanceRepo, err2 = postgresRepo.NewBalanceRepo(ctx, pgconn)
	if err2 != nil {
		err = errors.Wrap(err2, "couldn't make balance repo")
		return
	}

	statsRepo, err2 = postgresRepo.NewMarketStatsRepo(ctx, pgconn)
	if err2 != nil {
		err = errors.Wrap(err2, "couldn't make stats repo")
		return
	}

	exchangeRepo, err2 = postgresRepo.NewExchangeRepo(ctx, pgconn)
	if err2 != nil {
		err = errors.Wrap(err2, "couldn't make exchange repo")
		return
	}

	return

}

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
		marketsRepo, balanceRepo, statsRepo, exchangeRepo, err := getRepos(config.Database)
		if err != nil {
			logger.Fatal("couldn't get repos", zap.Error(err))
		}
		if marketsRepo == nil || balanceRepo == nil || statsRepo == nil {
			logger.Fatal("one of repos is nil", zap.Error(err))
		}
		logger.Info("all repos initialized")

		wallexClient, err := wallex.NewWallexClient(config.Exchange.Wallex, marketsRepo)
		if err != nil {
			logger.Fatal("couldn't connect to wallex", zap.Error(err))
		}

		ramzinexClient, err := ramzinex.NewRamzinexClient(config.Exchange.Ramzinex, marketsRepo)
		if err != nil {
			logger.Fatal("couldn't connect to ramzinex", zap.Error(err))
		}
		logger.Info("all exchanges connected")

		tgbot, err := telegrambot.NewTelegramBot(config.Telegram, pkg.GetLogger("telegram_bot").With(zap.String("package", "telegram bot")))
		if err != nil {
			logger.Panic("couldn't initialize telegram bot", zap.Error(err))
		}
		logger.Info("telegram bot initialized")

		priceManager := core.NewPriceManager(statsRepo, logger.With(zap.String("module", "price manager")))
		parooCore, err := core.NewParooCore(tgbot, []exchange.Exchange{wallexClient, ramzinexClient}, balanceRepo, marketsRepo,
			statsRepo, priceManager, exchangeRepo)
		if err != nil {
			logger.Panic("couldn't initialize core", zap.Error(err))
		}

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

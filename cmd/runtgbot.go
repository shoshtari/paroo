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

func GetRepos(config configs.SectionDatabase) (
	marketRepo repositories.MarketRepo,
	balanceRepo repositories.BalanceRepo,
	statsRepo repositories.MarketStatsRepo,
	err error) {

	if config.Provider == "" {
		if config.Postgres.Host != "" {
			config.Provider = "postgres"
		} else {
			config.Provider = "sqlite"
		}
	}

	ctx := context.Background()
	switch config.Provider {
	case "postgres":
		pkg.GetLogger().Info("using postgres for db")
		pgconn, err2 := postgresRepo.ConnectPostgres(ctx, config.Postgres)
		if err2 != nil {
			err = errors.Wrap(err2, "couldn't connect to postgres")
			return
		}

		marketRepo, err2 = postgresRepo.NewMarketRepo(pgconn, ctx)
		if err2 != nil {
			err = errors.Wrap(err, "couldn't make markets repo")
			return
		}

		balanceRepo, err2 = postgresRepo.NewBalanceRepo(pgconn, ctx)
		if err2 != nil {
			err = errors.Wrap(err, "couldn't make balance repo")
			return
		}

		statsRepo, err2 = postgresRepo.NewMarketStatsRepo(pgconn, ctx)
		if err2 != nil {
			err = errors.Wrap(err, "couldn't make stats repo")
			return
		}
		return

	case "sqlite":
		pkg.GetLogger().Info("using sqlite for db")

		db, err2 := sqliteRepo.Connect(config.Sqlite)
		if err2 != nil {
			err = err2
			return
		}

		marketRepo, err2 = sqliteRepo.NewMarketRepo(ctx, db)
		if err != nil {
			err = errors.WithMessage(err2, "couldn't initialize market repo")
		}

		balanceRepo, err2 = sqliteRepo.NewBalanceRepo(ctx, db)
		if err != nil {
			err = errors.WithMessage(err2, "couldn't initialize balance  repo")
		}

		statsRepo, err2 = sqliteRepo.NewMarketStatsRepo(ctx, db)
		if err != nil {
			err = errors.WithMessage(err2, "couldn't initialize stats  repo")
		}
		return
	}
	return

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
		marketsRepo, balanceRepo, statsRepo, err := GetRepos(config.Database)
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

		parooCore := core.NewParooCode(tgbot, wallexClient, balanceRepo, statsRepo)

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

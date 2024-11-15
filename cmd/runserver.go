package cmd

import (
	"log"

	"github.com/shoshtari/paroo/internal"
	"github.com/shoshtari/paroo/internal/configs"
	"github.com/spf13/cobra"
)

// runserverCmd represents the runserver command
var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "run the server",
	Long:  "run the server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configs.GetConfig()
		if err != nil {
			log.Fatal("couldn't get config, err: ", err)
		}
		logger, err := internal.GetLogger(config.Log)
		if err != nil {
			log.Fatal("couldn't initialize logger, err: ", err)
		}

		internal.RunServer(config.HTTPServer, logger)
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)
}

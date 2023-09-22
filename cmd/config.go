package cmd

import (
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "root",
		Short: "root",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func readConfig(cmd *cobra.Command, args []string) {
	if err := envconfig.Process("", &cfg); err != nil {
		log.Err(err).Msg("")
		return
	}

	l, err := zerolog.ParseLevel(strings.ToLower(cfg.Log.Level))
	if err != nil {
		log.Err(err).Msg("")
		return
	}
	zerolog.SetGlobalLevel(l)

	if strings.ToLower(cfg.Log.Format) == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/negativeOne1/petcam/internal/config"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:              "worker",
	Short:            "worker",
	PersistentPreRun: preRun,
}

func init() {
	var err error
	cfg, err = config.Load("")
	if err != nil {
		log.Fatal().Err(err).Msg("can not load config")
	}

	err = pdftools.Initialize(cfg.App.PDFToolsLicenseKey)
	if err != nil {
		log.Fatal().Err(err).Msg("could not initialize pdftools sdk: %s")
	}
	defer pdftools.UnInitialize()

	cCmd := worker.NewRPCCobraCommand("rpc-compress", compress.ProcessFn[compress.Task])
	cCmd.Run = runMiddleware(cCmd.Run)
	rootCmd.AddCommand(cCmd)
	worker.RegisterTypes(
		&compress.Task{},
		&compress.Result{},
	)
}

func Execute() error {
	return rootCmd.Execute()
}

func preRun(cmd *cobra.Command, args []string) {
	var err error

	l, err := zerolog.ParseLevel(strings.ToLower(cfg.Log.Level))
	if err != nil {
		log.Fatal().Err(err).Msg("can not parse log level")
		return
	}
	zerolog.SetGlobalLevel(l)

	if strings.ToLower(cfg.Log.Format) == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

var runMiddleware = func(next func(*cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if cfg.DataDog.InstrumentationEnabled {
			tracer.Start(tracer.WithRuntimeMetrics())
			defer tracer.Stop()
		}

		go func() {
			time.Sleep(30 * time.Minute)
			log.Info().Msg("worker timelimit exceeded, shutting down")
			os.Exit(0)
		}()

		next(cmd, args)
	}
}

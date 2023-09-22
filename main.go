package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Err(err).Msg("")
	}
}

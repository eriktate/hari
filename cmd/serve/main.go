package main

import (
	"os"

	"github.com/eriktate/hari/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func run() error {
	log.Info().Msg("starting hari...")

	return http.Serve(http.Config{Addr: "0.0.0.0:9001"})
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if err := run(); err != nil {
		log.Error().Err(err).Msg("hari server failure")
	}
}

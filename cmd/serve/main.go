package main

import (
	"fmt"
	"os"

	"github.com/eriktate/hari/http"
	"github.com/eriktate/hari/pg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func run() error {
	log.Info().Msg("starting hari...")

	pg, err := pg.New("postgres://hari:hari@localhost:5432/hari?sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to init connection with postgres: %w", err)
	}

	cfg := http.Config{
		Addr:           "0.0.0.0:9001",
		WebhookService: pg,
	}

	return http.Serve(cfg)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if err := run(); err != nil {
		log.Error().Err(err).Msg("hari server failure")
	}
}

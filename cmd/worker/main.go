package main

import (
	"context"
	"flag"

	"github.com/Iwoooooods/fs-upload-go/core/database"
	"github.com/Iwoooooods/fs-upload-go/pkg/config"
	"github.com/rs/zerolog/log"
)

func main() {
	var envDir string
	flag.StringVar(&envDir, "env-dir", "./dev.env", "path to the environment directory")
	flag.Parse()

	log.Info().Str("env-dir", envDir).Msg("reading config from: " + envDir)
	cfg := config.Load(envDir)

	conn := database.ConnectSqlite(cfg.DbName)
	ctx := context.Background()
	if err := conn.Connect(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to sqlite")
	}

	log.Info().Msg("connected to sqlite")
}

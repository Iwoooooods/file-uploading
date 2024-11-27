package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Iwoooooods/fs-upload-go/api"
	"github.com/Iwoooooods/fs-upload-go/internal/config"
	"github.com/Iwoooooods/fs-upload-go/internal/database"
	"github.com/rs/zerolog/log"
)

const (
	ACCESS_TOKEN_HEADER = "X-Access-Token"
	STREAM_TOKEN_HEADER = "X-Stream-Token"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.Parse()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, ACCESS_TOKEN_HEADER, STREAM_TOKEN_HEADER},
		ExposeHeaders: []string{echo.HeaderContentLength, echo.HeaderContentDisposition, echo.HeaderContentEncoding},
	}))

	cfg := config.Load(".env")
	log.Info().Msg("reading config from: " + ".env")

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to create database")
	}
	defer db.Close()

	router := e.Group("api")
	apiHandler := api.NewHandler(cfg, db)
	apiHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: e,
	}

	go func() {
		log.Info().Str("port", port).Msg("server started at: " + fmt.Sprintf(":%s", port))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server startup failed")
		}
	}()

	appCtx := context.Background()
	// listen for os interrupt signals
	ctx, cancel := signal.NotifyContext(appCtx, os.Interrupt)
	defer cancel()

	// block until user interrupts the program (ctrl+c)
	<-ctx.Done()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("error during server shutdown")
	}
}

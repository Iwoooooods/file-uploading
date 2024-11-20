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

	"github.com/rs/zerolog/log"
)

const ACCESS_TOKEN_HEADER = "X-Access-Token"

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func main() {
	var port string
	flag.StringVar(&port, "port", "9191", "port to listen on")
	flag.Parse()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, ACCESS_TOKEN_HEADER},
		ExposeHeaders: []string{echo.HeaderContentLength, echo.HeaderContentDisposition, echo.HeaderContentEncoding},
	}))

	router := e.Group("api")

	router.GET("/ping", hello)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: e,
	}

	go func() {
		log.Info().Str("port", port).Msg("server started at")

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

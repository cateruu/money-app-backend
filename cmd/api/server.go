package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cateruu/money-app-backend/internal/data"
	"github.com/cateruu/money-app-backend/pkg/logger"
)

func (app *app) serve(models *data.Models) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(models),
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Log.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Log.Info("Shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.wg.Wait()
		shutdownError <- nil
	}()

	logger.Log.Info("Starting server", "addr", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	logger.Log.Info("Stopped server", "addr", srv.Addr)
	return nil

}

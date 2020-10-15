package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/starrybarry/schedule/internal/amqplb"

	"github.com/starrybarry/schedule/internal/scheduler"

	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/starrybarry/schedule/cmd/scheduler/cfg"
	"github.com/starrybarry/schedule/cmd/scheduler/handler"
	"go.uber.org/zap"
)

func main() {
	rootCtx, cancel := context.WithCancel(context.Background())

	streamErr := make(chan error, 0)

	log := zap.NewExample()

	config, err := cfg.NewConfig()
	if err != nil {
		log.Fatal("new config", zap.Error(err))
	}

	clientAMQP := amqplb.NewClient(config.Rabbit.DSN, config.Rabbit.Heartbeat)
	if err = clientAMQP.Connect(); err != nil {
		log.Fatal("connect rabbit", zap.Error(err), zap.String("rabbit_dsn", config.Rabbit.DSN))
	}

	amqpClose := func() {
		if errr := clientAMQP.Close(); errr != nil {
			log.Error("close client amqp", zap.Error(errr))
		}
	}

	defer amqpClose()

	pgxPool, err := pgxpool.Connect(rootCtx, config.Postgre.URL)
	if err != nil {
		log.Fatal("postgreSQL connection error: ", zap.Error(err))
	}

	defer pgxPool.Close()

	log.Info("postgre created pool!", zap.String("url", config.Postgre.URL))

	schStorage := scheduler.NewTaskStorage(pgxPool, log)
	schService := scheduler.NewService(schStorage)
	handlerHTTP := handler.NewHttpHandler(schService, log)

	if config.DebugMode {
		handlerHTTP = cors.AllowAll().Handler(handlerHTTP)
		log.Info("cors enabled!")
	}

	srv := &http.Server{
		Addr:           config.ServeAddr,
		Handler:        handlerHTTP,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 64 * 1024,
	}

	go func() {
		log.Info("Start server...",
			zap.String("service_name", config.ServiceName), zap.Any("address", srv.Addr))

		if err := srv.ListenAndServe(); err != nil {
			streamErr <- errors.Wrap(err, "failed to listen and serve server")
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-streamErr:
		log.Error("error in errCh", zap.Error(err))

		if errs := srv.Shutdown(rootCtx); errs != nil {
			log.Error("failed tp shutdown  serve", zap.Error(errs))
		}

		cancel()

		log.Info("Server shutdown")
	case <-sig:
		if err := srv.Shutdown(rootCtx); err != nil {
			log.Error("failed tp shutdown  serve", zap.Error(err))
		}

		cancel()

		log.Info("Server shutdown")
	}
}

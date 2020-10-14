package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/starrybarry/schedule/cmd/migration/cfg"
	"go.uber.org/zap"
)

func main() {
	log := zap.NewExample()

	config, err := cfg.NewConfig()
	if err != nil {
		log.Fatal("new config", zap.Error(err))
	}

	log.Info("config", zap.Any("postgre", config.Postgre))

	m, err := migrate.New(config.Postgre.MigrationsPath, config.Postgre.URL)
	if err != nil {
		log.Fatal("new migrate", zap.Error(err))
	}

	if err = m.Up(); err == nil {
		log.Info("migrate up success!")
	}

	if err != migrate.ErrNoChange {
		log.Fatal("up migration", zap.Error(err))
	}
}

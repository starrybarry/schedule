package cfg

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const (
	PostgreURL     = "postgre-url"
	MigrationsPath = "migrations-path"
)

func ParseAllFlags() error {
	parseCommandLineFlags()
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return errors.Wrap(err, "failed to get flags from command line")
	}

	return nil
}

func parseCommandLineFlags() {
	pflag.String(PostgreURL, "postgres://scheduler:0000@localhost:5432/starry?sslmode=disable", "address for postgre, dsn")
	pflag.String(MigrationsPath, "file://./cmd/migration/migrations", "path to migrations")
	pflag.Parse()
}

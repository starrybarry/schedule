package cfg

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const (
	DebugMode = "debug"
	ServeAddr = "serve-addr"
	ServeName = "serve-name"

	//
	PostgreURL = "postgre-url"
	//
	RabbitDSN       = "rabbitmq-dsn"
	RabbitExchange  = "rabbit-exchange"
	RabbitHeartbeat = "rabbit-heartbeat"
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
	pflag.Bool(DebugMode, false, "is debug?")
	pflag.String(ServeAddr, ":8080", fmt.Sprintf(" serve addr\n env: %s\n", ServeAddr))
	pflag.String(ServeName, "scheduler", fmt.Sprintf(" serve name\n env: %s\n", ServeAddr))
	//

	pflag.String(RabbitDSN, "amqp://guest:guest@127.0.0.1:5672", "address for rabbit, dsn")
	pflag.String(RabbitExchange, "parser-data-store.parser-events", "exchange for parser events")
	pflag.Duration(RabbitHeartbeat, 10*time.Second, "rabbit heartbeat")
	//

	pflag.String(PostgreURL, "postgres://scheduler:0000@localhost:5432/scheduler", "address for postgre, dsn")

	pflag.Parse()
}

package shared

import (
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

// Logger is an instance of the shared logger tool.
var Logger *logging.Logger

const (
	format = "%{color}%{time:2006.01.02 15:04:05} " +
		"%{id:04x} %{level:.4s}%{color:reset} " +
		"[%{module}] %{color:bold}%{shortfunc}%{color:reset} -> %{message}"
)

func initLogger() {
	var (
		envLevel = viper.GetString("logging")
		chaincodeName = viper.GetString("name")
	)

	Logger = logging.MustGetLogger(chaincodeName)

	backend := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(format),
	)

	level, err := logging.LogLevel(envLevel); if err != nil {
		Logger.Warningf("failed to parse '%s' as a log level", envLevel)
		level = logging.INFO
	}

	logging.SetBackend(backend)
	logging.SetLevel(level, chaincodeName)
}

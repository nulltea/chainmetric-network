package shared

import (
	"os"

	"github.com/op/go-logging"
)

var (
	// Logger is the shared logger instance.
	Logger = logging.MustGetLogger(ChaincodeName)
)

const (
	format = "%{color}%{time:2006.01.02 15:04:05} %{id:04x} %{level:.4s}%{color:reset} [%{module}] %{color:bold}%{shortfunc}%{color:reset} -> %{message}"
)

func initLogger() {

	backend := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(format),
	)

	level, err := logging.LogLevel(LogLevel); if err != nil {
		Logger.Warningf("failed to parse '%s' as a log level", LogLevel)
		level = logging.INFO
	}

	logging.SetBackend(backend)
	logging.SetLevel(level, ChaincodeName)
}

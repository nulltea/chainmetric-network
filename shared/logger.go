package shared

import (
	"os"

	"github.com/op/go-logging"
)

const (
	format = "%{color}%{time:2006.01.02 15:04:05} %{id:04x} %{level:.4s}%{color:reset} [%{module}] %{color:bold}%{shortfunc}%{color:reset} -> %{message}"
)

var (
	Logger = logging.MustGetLogger(ChaincodeName)
)

func InitLogger() {
	backend := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(format))
	logging.SetBackend(backend)

	level, err := logging.LogLevel(LogLevel); if err != nil {
		level = logging.WARNING
	}
	logging.SetLevel(level, ChaincodeName)
}

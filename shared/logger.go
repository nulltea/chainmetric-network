package shared

import "github.com/op/go-logging"

var (
	Logger = logging.MustGetLogger(ChaincodeName)
	format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000 MST} %{id:04x} %{level:.4s}%{color:reset} [%{module}] %{color:bold}%{shortfunc}%{color:reset} -> %{message}`,
	)
)

func InitLogger() {
	logging.SetFormatter(format)
}

package config

import (
	"os"

	"github.com/op/go-logging"
)

var backend = logging.NewLogBackend(os.Stderr, "", 0)
var format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{level:.3s}][%{module}] %{message}`)

// MakeLogger creates a regular logger
func MakeLogger(module string) *logging.Logger {
	backendFormatted := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatted)
	l := &logging.Logger{Module: module}
	l.SetBackend(logging.AddModuleLevel(backendFormatted))
	return l
}

// ConfigureLogger sets up the logger and it's backends
func ConfigureLogger() {
	backendFormatted := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatted)
}

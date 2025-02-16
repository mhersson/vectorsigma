package statemachine

import (
	"log/slog"
)

type Context struct {
	Log *slog.Logger
}

type ExtendedState struct {
	Error string
}

package fsm

import (
	"log/slog"
)

// A struct that holds the items needed for the actions to do their work.
// Things like client libraries and loggers, go here.
type Context struct {
	Logger *slog.Logger // Do NOT delete this!
}

// A struct that holds the "extended state" of the state machine, including data
// being fetched and read. This should only be modified by actions, while guards
// should only read the extended state to assess their value.
type ExtendedState struct {
	Error error
}

package main

import (
	"log/slog"
	"initializeNewModule/pkg/fsm"
)

func main() {
	machine := fsm.New()
	runner := fsm.Runner{
		StateMachine: *machine,
		Context: &fsm.Context{
			// Log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
			Log: slog.Default(),
		},
		State:   &fsm.ExtendedState{},
		Actions: fsm.AllActions,
		Guards:  fsm.AllGuards,
	}

	runner.Run()
}

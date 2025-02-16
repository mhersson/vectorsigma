package statemachine_test

import (
	"errors"
	"testing"

	"log/slog"

	"github.com/mhersson/vectorsigma/internal/statemachine"
)

func TestFSM_IsErrorGuard(t *testing.T) {
	type fields struct {
		logger        *slog.Logger
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "IsError", fields: fields{ExtendedState: &statemachine.ExtendedState{Error: errors.New("error")}}, want: true},
		{name: "NoError", fields: fields{ExtendedState: &statemachine.ExtendedState{Error: nil}}, want: false},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.FSM{
				Logger:        tt.fields.logger,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.IsErrorGuard(); got != tt.want {
				t.Errorf("FSM.IsErrorGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_IsMarkdownGuard(t *testing.T) {
	type fields struct {
		logger        *slog.Logger
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "Is Markdown", fields: fields{ExtendedState: &statemachine.ExtendedState{Input: "input.md"}}, want: true},
		{name: "Is Plantuml", fields: fields{ExtendedState: &statemachine.ExtendedState{Input: "input.plantuml"}}, want: false},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.FSM{
				Logger:        tt.fields.logger,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.IsMarkdownGuard(); got != tt.want {
				t.Errorf("FSM.IsMarkdownGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_IsStandaloneModuleGuard(t *testing.T) {
	type fields struct {
		logger        *slog.Logger
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "Is Standalone Module", fields: fields{ExtendedState: &statemachine.ExtendedState{Init: true}}, want: true},
		{name: "Is only a package", fields: fields{ExtendedState: &statemachine.ExtendedState{Init: false}}, want: false},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.FSM{
				Logger:        tt.fields.logger,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.IsStandaloneModuleGuard(); got != tt.want {
				t.Errorf("FSM.IsStandaloneModuleGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

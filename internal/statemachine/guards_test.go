package statemachine_test

import (
	"errors"
	"testing"

	"github.com/mhersson/vectorsigma/internal/statemachine"
)

func TestVectorSigma_IsErrorGuard(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
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
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.IsErrorGuard(); got != tt.want {
				t.Errorf("VectorSigma.IsErrorGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorSigma_IsMarkdownGuard(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
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
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.IsMarkdownGuard(); got != tt.want {
				t.Errorf("VectorSigma.IsMarkdownGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorSigma_IsStandaloneModuleGuard(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
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
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.IsStandaloneModuleGuard(); got != tt.want {
				t.Errorf("VectorSigma.IsStandaloneModuleGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorSigma_PackageExistsGuard(t *testing.T) {
	type fields struct {
		context       *statemachine.Context
		currentState  statemachine.StateName
		stateConfigs  map[statemachine.StateName]statemachine.StateConfig
		ExtendedState *statemachine.ExtendedState
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "Package exists", fields: fields{ExtendedState: &statemachine.ExtendedState{PackageExits: true}}, want: true},
		{name: "Package does not exist", fields: fields{ExtendedState: &statemachine.ExtendedState{PackageExits: false}}, want: false},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &statemachine.VectorSigma{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.PackageExistsGuard(); got != tt.want {
				t.Errorf("VectorSigma.PackageExistsGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

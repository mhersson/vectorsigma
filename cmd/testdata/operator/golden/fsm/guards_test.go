package fsm_test

import (
	"operator/output/fsm"
	"testing"
)

// +vectorsigma:guard:IsError
func TestTestreconcileloop_IsErrorGuard(t *testing.T) {
	type fields struct {
		context       *fsm.Context
		currentState  fsm.StateName
		stateConfigs  map[fsm.StateName]fsm.StateConfig
		ExtendedState *fsm.ExtendedState
	}
	type args struct {
		params []string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &fsm.Testreconcileloop{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.IsErrorGuard(tt.args.params...); got != tt.want {
				t.Errorf("Testreconcileloop.IsErrorGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

// +vectorsigma:guard:NotFound
func TestTestreconcileloop_NotFoundGuard(t *testing.T) {
	type fields struct {
		context       *fsm.Context
		currentState  fsm.StateName
		stateConfigs  map[fsm.StateName]fsm.StateConfig
		ExtendedState *fsm.ExtendedState
	}
	type args struct {
		params []string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &fsm.Testreconcileloop{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.NotFoundGuard(tt.args.params...); got != tt.want {
				t.Errorf("Testreconcileloop.NotFoundGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

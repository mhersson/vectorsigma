package fsm_test

import (
	"testing"

	"package/fsm"
)

func TestTrafficLight_IsErrorGuard(t *testing.T) {
	type fields struct {
		context       *fsm.Context
		currentState  fsm.StateName
		stateConfigs  map[fsm.StateName]fsm.StateConfig
		ExtendedState *fsm.ExtendedState
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsm := &fsm.TrafficLight{
				Context:       tt.fields.context,
				CurrentState:  tt.fields.currentState,
				StateConfigs:  tt.fields.stateConfigs,
				ExtendedState: tt.fields.ExtendedState,
			}
			if got := fsm.IsErrorGuard(); got != tt.want {
				t.Errorf("TrafficLight.IsErrorGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

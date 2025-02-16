package integrationtest_test

import (
	"testing"

	"new_module/integrationtest"
)

func TestTrafficLight_IsErrorGuard(t *testing.T) {
	type fields struct {
		context       *integrationtest.Context
		currentState  integrationtest.StateName
		stateConfigs  map[integrationtest.StateName]integrationtest.StateConfig
		ExtendedState *integrationtest.ExtendedState
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
			fsm := &integrationtest.TrafficLight{
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
